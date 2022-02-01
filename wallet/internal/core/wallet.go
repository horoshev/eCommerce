package core

import (
	"context"
	"eCommerce/wallet/internal/models"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type WalletService interface {
	CreateNewWallet(user *models.User) (*models.Wallet, error)
	PayOrder(order *models.Order) (*models.Transaction, error)
	CancelOrder(order *models.Order) (*models.Transaction, error)
}

type WalletController struct {
	ctx     context.Context
	log     *zap.SugaredLogger
	wallets *mongo.Collection
	//transactions *mongo.Collection
	producer *kafka.Writer
}

func NewWalletController(ctx context.Context, log *zap.SugaredLogger, db *mongo.Database, producer *kafka.Writer) *WalletController {
	wallet := new(WalletController)

	wallet.ctx = ctx
	wallet.log = log
	wallet.wallets = db.Collection(`wallets`)
	wallet.producer = producer

	return wallet
}

func (w *WalletController) CreateNewWallet(user *models.User) (*models.Wallet, error) {
	wallet := new(models.Wallet)
	wallet.UserId = user.Id
	wallet.UserName = user.Name
	wallet.Balance = 100
	wallet.Transactions = []models.Transaction{
		{
			OrderId:   user.Id,
			Status:    models.TransactionActive,
			Note:      "New Customer Bonus",
			Amount:    100,
			Balance:   100,
			Timestamp: time.Now().UTC(),
		},
	}

	single, err := w.wallets.InsertOne(context.Background(), wallet)
	if err != nil {
		return nil, err
	}

	wallet.Id = single.InsertedID.(primitive.ObjectID)

	value, err := json.Marshal(wallet)
	if err != nil {
		w.log.Error(err)
	}

	err = w.producer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(user.Id.Hex()),
		Value: value,
		Topic: models.WalletResponseTopic,
	})
	if err != nil {
		w.log.Error(err)
	}

	return wallet, nil
}

func (w *WalletController) PayOrder(key primitive.ObjectID, order *models.Order) (*models.Transaction, error) {
	response := new(Response)
	defer func() {
		message, err := response.Marshal()
		if err != nil {
			w.log.Error(err)
			message = []byte(`marshaling error`)
		}
		err = w.producer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(key.Hex()),
			Value: message,
			Topic: models.WalletPayOrderResponseTopic,
			Headers: []kafka.Header{
				{Key: `status`, Value: response.StatusHeader()},
				{Key: `message`, Value: []byte(response.Message)},
			},
		})
		if err != nil {
			w.log.Error(err)
		}
	}()

	if order == nil {
		err := errors.New(`argument 'order' is nil`)
		response = NewError(err, order)
		return nil, err
	}

	wallet, err := w.findUserWallet(order.UserId)
	if err != nil {
		response = NewError(err, order)
		return nil, err
	}

	if wallet.Balance < order.Amount {
		err = errors.New("there are not enough funds on the wallet")
		response = NewError(err, order)
		return nil, err
	}

	transaction := NewOrderPayment(order.Id, order.Amount, wallet.Balance)
	filter := bson.D{{"_id", wallet.Id}}
	update := bson.D{
		{"$push", bson.M{`transactions`: transaction}},
		{"$set", bson.M{`balance`: transaction.Balance}},
	}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	single := w.wallets.FindOneAndUpdate(context.Background(), filter, update, option)
	if err = single.Err(); err != nil {
		response = NewError(err, order)
		return nil, err
	}

	response = NewSuccess(`payment successful`, order)

	return transaction, nil
}

func NewOrderPayment(orderId primitive.ObjectID, amount, balance float64) *models.Transaction {
	transaction := new(models.Transaction)
	transaction.OrderId = orderId
	transaction.Status = models.TransactionActive
	transaction.Note = "Order payment"
	transaction.Amount = -amount
	transaction.Balance = balance - amount
	transaction.Timestamp = time.Now().UTC()

	return transaction
}

func (w *WalletController) CancelOrder(order *models.Order) (*models.Transaction, error) {
	response := new(Response)
	defer func() {
		message, err := response.Marshal()
		if err != nil {
			w.log.Error(err)
			message = []byte(`marshaling error`)
		}
		err = w.producer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(order.Id.Hex()),
			Value: message,
			Topic: models.WalletCancelOrderResponseTopic,
			Headers: []kafka.Header{
				{Key: `status`, Value: response.StatusHeader()},
				{Key: `message`, Value: []byte(response.Message)},
			},
		})
		if err != nil {
			w.log.Error(err)
		}
	}()

	wallet, err := w.findUserWallet(order.Id)
	if err != nil {
		response = NewError(err, order)
		return nil, err
	}

	transaction, err := w.findTransaction(order.Id)
	if err != nil {
		response = NewError(err, order)
		return nil, err
	}

	if wallet.Balance-transaction.Amount < 0 {
		err = errors.New("there are not enough funds on the wallet")
		response = NewError(err, order)
		return nil, err
	}

	revert := RevertOrderPayment(transaction, wallet.Balance)
	filter := bson.D{
		{"_id", order.UserId},
		{"transactions.order_id", order.Id},
		{"transactions.status", models.TransactionActive},
	}
	update := bson.D{
		{"$push", bson.M{`transactions`: revert}},
		{"$set", bson.M{`transactions.$.status`: models.TransactionCancelled}},
	}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	single := w.wallets.FindOneAndUpdate(context.Background(), filter, update, option)
	if err = single.Err(); err != nil {
		response = NewError(err, order)
		return nil, err
	}

	response = NewSuccess("canceled order payment", order)

	return revert, nil
}

func RevertOrderPayment(transaction *models.Transaction, balance float64) *models.Transaction {
	revert := new(models.Transaction)
	revert.Status = models.TransactionActive
	revert.Note = `Cancel of order payment ` + transaction.Id.Hex()
	revert.Amount = -transaction.Amount
	revert.Balance = balance - transaction.Amount
	revert.Timestamp = time.Now().UTC()

	return revert
}

// findUserWallet fetch wallet data from database.
func (w *WalletController) findUserWallet(id primitive.ObjectID) (*models.Wallet, error) {
	single := w.wallets.FindOne(context.Background(), bson.M{"user_id": id})
	if err := single.Err(); err != nil {
		return nil, err
	}

	wallet := new(models.Wallet)
	err := single.Decode(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// findTransaction fetch transaction data from database.
func (w *WalletController) findTransaction(orderId primitive.ObjectID) (*models.Transaction, error) {
	single := w.wallets.FindOne(context.Background(), bson.M{"order_id": orderId})
	if err := single.Err(); err != nil {
		return nil, err
	}

	transaction := new(models.Transaction)
	if err := single.Decode(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// updateTransactionStatus find a transaction and updates it status.
//func (w *WalletController) updateTransactionStatus(id primitive.ObjectID, status models.TransactionStatus) error {
//	filter := bson.D{{"_id", id}}
//	update := bson.D{{"$set", bson.D{{"status", status}}}}
//	one := w.transactions.FindOneAndUpdate(context.Background(), filter, update)
//
//	if err := one.Err(); err != nil {
//		//c.log.Fatal(err)
//		return err
//	}
//
//	return nil
//}
