package core

import (
	"context"
	"eCommerce/storage/internal/models"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
)

const (
	StorageReserveOrderResponseTopic = `storage-reserve-order-response`
	StorageCancelOrderResponseTopic  = `storage-cancel-order-response`
)

var (
	wc        = writeconcern.New(writeconcern.WMajority())
	rc        = readconcern.Snapshot()
	txOptions = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
)

type StorageService interface {
	ReserveOrder(order *models.Order) error
	CancelOrder(order *models.Order) error
}

type Storage struct {
	log          *zap.SugaredLogger
	products     *mongo.Collection
	reservations *mongo.Collection
	producer     *kafka.Writer
}

func NewStorage(log *zap.SugaredLogger, database *mongo.Database, producer *kafka.Writer) *Storage {
	storage := new(Storage)
	storage.log = log
	storage.products = database.Collection(`products`)
	storage.reservations = database.Collection(`reservations`)
	storage.producer = producer

	return storage
}

func (s Storage) ReserveOrder(order *models.Order) error {
	response := new(Response)
	defer func() {
		message, err := response.Marshal()
		if err != nil {
			message = []byte(`marshaling error`)
		}

		err = s.producer.WriteMessages(context.TODO(), kafka.Message{
			Key:   []byte(order.Id.Hex()),
			Value: message,
			Topic: StorageReserveOrderResponseTopic,
			Headers: []kafka.Header{
				{Key: `status`, Value: response.StatusHeader()},
				{Key: `message`, Value: []byte(response.Message)},
			},
		})
		if err != nil {
			s.log.Error(err)
		}
	}()

	dbSession, err := s.products.Database().Client().StartSession()
	if err != nil {
		response = NewError(err, order)
		return err
	}
	defer dbSession.EndSession(context.TODO())

	err = mongo.WithSession(context.Background(), dbSession, s.ReserveOrderTx(dbSession, order))
	if err != nil {
		response = NewError(err, order)
		return err
	}

	response = NewSuccess("reserved order", order)

	return nil
}

func (s Storage) CancelOrder(order *models.Order) error {
	response := new(Response)
	defer func() {
		message, err := response.Marshal()
		if err != nil {
			message = []byte(`marshaling error`)
		}
		err = s.producer.WriteMessages(context.TODO(), kafka.Message{
			Key:   []byte(order.Id.Hex()),
			Value: message,
			Topic: StorageCancelOrderResponseTopic,
			Headers: []kafka.Header{
				{Key: `status`, Value: response.StatusHeader()},
				{Key: `message`, Value: []byte(response.Message)},
			},
		})
		if err != nil {
			s.log.Error(err)
		}
	}()

	dbSession, err := s.products.Database().Client().StartSession()
	if err != nil {
		response = NewError(err, order)
		return err
	}
	defer dbSession.EndSession(context.TODO())

	err = mongo.WithSession(context.Background(), dbSession, s.CancelOrderTx(dbSession, order))
	if err != nil {
		response = NewError(err, order)
		return err
	}

	response = NewSuccess("canceled order", order)

	return nil
}

func (s Storage) ReserveOrderTx(session mongo.Session, order *models.Order) func(ctx mongo.SessionContext) error {
	return func(ctx mongo.SessionContext) error {
		if err := session.StartTransaction(txOptions); err != nil {
			s.log.Error(err)
			return err
		}

		amount, err := s.OrderAmount(ctx, order)
		if err != nil {
			_ = session.AbortTransaction(ctx)
			s.log.Error(err)
			return err
		}

		order.Amount = amount
		if err = s.DeductProducts(ctx, order); err != nil {
			_ = session.AbortTransaction(ctx)
			s.log.Error(err)
			return err
		}

		reservation := models.OrderToReservation(order)
		_, err = s.reservations.InsertOne(context.Background(), reservation)
		if err != nil {
			s.log.Error(err)
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			s.log.Error(err)
			return err
		}

		return nil
	}
}

func (s Storage) CancelOrderTx(session mongo.Session, order *models.Order) func(ctx mongo.SessionContext) error {
	return func(ctx mongo.SessionContext) error {
		if err := session.StartTransaction(txOptions); err != nil {
			return err
		}

		reservation, err := s.IsReserved(ctx, order)
		if err != nil {
			_ = session.AbortTransaction(ctx)
			return errors.New("out of order")
		}

		if err = s.ReturnProducts(ctx, reservation); err != nil {
			_ = session.AbortTransaction(ctx)
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		return nil
	}
}

func (s Storage) OrderAmount(ctx mongo.SessionContext, order *models.Order) (amount float64, err error) {
	for _, p := range order.Items {
		filter := bson.D{
			{"name", p.Name},
			{"quantity", bson.D{{"$gte", p.Quantity}}},
		}

		single := s.products.FindOne(ctx, filter)
		if single.Err() != nil {
			return 0, errors.New(`one of the products out of stock`)
		}

		product := new(models.Product)
		err = single.Decode(product)
		if err != nil {
			return 0, err
		}

		amount += product.Cost * float64(p.Quantity)
	}

	return amount, nil
}

func (s Storage) IsInStock(ctx mongo.SessionContext, order *models.Order) bool {
	for _, p := range order.Items {
		filter := bson.D{
			{"name", p.Name},
			{"quantity", bson.D{{"$gte", p.Quantity}}},
		}

		single := s.products.FindOne(ctx, filter)
		if single.Err() != nil {
			return false
		}
	}

	return true
}

func (s Storage) IsReserved(ctx mongo.SessionContext, order *models.Order) (*models.OrderReservation, error) {
	filter := bson.D{
		{"order_id", order.Id},
		{"status", models.Success},
	}
	single := s.reservations.FindOne(ctx, filter)
	if err := single.Err(); err != nil {
		return nil, err
	}

	reservation := new(models.OrderReservation)
	if err := single.Decode(reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s Storage) DeductProducts(ctx mongo.SessionContext, order *models.Order) error {
	for _, p := range order.Items {
		filter := bson.D{{"name", p.Name}}
		update := bson.D{{"$inc", bson.D{{"quantity", -p.Quantity}}}}
		single := s.products.FindOneAndUpdate(ctx, filter, update)
		if err := single.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) ReturnProducts(ctx mongo.SessionContext, order *models.OrderReservation) error {
	for _, p := range order.Products {
		filter := bson.D{{"name", p.ProductName}}
		update := bson.D{{"$inc", bson.D{{"quantity", p.Quantity}}}}
		single := s.products.FindOneAndUpdate(ctx, filter, update)
		if err := single.Err(); err != nil {
			return err
		}
	}

	filter := bson.D{{"_id", order.Id}}
	update := bson.D{{"$set", bson.D{{"status", models.Canceled}}}}
	s.reservations.FindOneAndUpdate(ctx, filter, update)

	return nil
}
