package core

import (
	"context"
	"eCommerce/registry/internal/consumers"
	"eCommerce/registry/internal/consumers/topics"
	"eCommerce/registry/internal/data"
	"eCommerce/registry/internal/models"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// OrderCoordinator is managing order flow.
type OrderCoordinator struct {
	log        *zap.SugaredLogger
	producer   *kafka.Writer
	consumers  *consumers.OrderConsumerSet
	repository data.RegistryRepository
}

func NewOrderCoordinator(log *zap.SugaredLogger, repository data.RegistryRepository, writer *kafka.Writer) *OrderCoordinator {
	oc := new(OrderCoordinator)
	oc.log = log
	oc.producer = writer
	oc.repository = repository
	oc.consumers = consumers.NewOrderConsumerSet(log, writer, repository)

	return oc
}

func (oc *OrderCoordinator) Run(baseConfig kafka.ReaderConfig) {
	oc.consumers.Run(baseConfig)
}

// NewOrder initialize new order saga and publish an event to kafka to reserve credit.
func (oc *OrderCoordinator) NewOrder(order *models.Order) {
	value, err := json.Marshal(order)
	if err != nil {
		oc.log.Error(err)
		return
	}

	err = oc.producer.WriteMessages(context.Background(), kafka.Message{
		Topic: topics.StorageReserveOrderTopic,
		Key:   []byte(order.Id.Hex()),
		Value: value,
	})
	if err != nil {
		oc.log.Error(err)
	}

	order, err = oc.repository.UpdateOrderStatus(order.Id, models.OrderReservationPending)
	if err != nil {
		oc.log.Error(err)
	}
}

/*
// OnOrderReserved processing products reservation result.
// When products successfully reserved - update order status to 'ORDER_RESERVED' and publish an event to kafka.
// When products reservation failed - update order status to 'Error' and publish an event to kafka to cancel credit reservation.
func (oc *OrderCoordinator) OnOrderReserved(order *models.Order) {

	// TODO:
	// ~ 1. receive a price from storage
	// + 2. send price to reserve in wallet
	// + 3. add update order method. also can update a amount of order
	// - 4. add repository

	order, err := oc.UpdateOrder(order.Id, bson.D{
		{"status", models.OrderReserved},
		{"amount", order.Amount},
	})
	if err != nil {
		oc.log.Error(err)
		return
	}

	value, err := json.Marshal(order)
	if err != nil {
		oc.log.Error(err)
		return
	}

	err = oc.Producer.WriteMessages(context.Background(), kafka.Message{
		Topic: WalletPayOrderTopic,
		Key:   []byte(order.Id.Hex()),
		Value: value,
	})
	if err != nil {
		oc.log.Error(err)
	}

	order, err = oc.UpdateOrderStatus(order.Id, models.OrderPaymentPending)
	if err != nil {
		oc.log.Error(err)
	}
}

// TODO
func (oc *OrderCoordinator) OnOrderReserveCanceled(order *models.Order) {
	order, err := oc.UpdateOrderStatus(order.Id, models.OrderReservationCanceled)
	if err != nil {
		oc.log.Error(err)
		return
	}

	order, err = oc.UpdateOrderStatus(order.Id, models.OrderCanceled)
	if err != nil {
		oc.log.Error(err)
		return
	}
}

// OnOrderPaid processing credit reservation result.
// When credit successfully reserved - update order status to 'CreditReserved' and publish an event to kafka to reserve products.
// When credit reservation failed - update order status to 'Error'.
func (oc *OrderCoordinator) OnOrderPaid(order *models.Order) {

	// TODO:
	//

	// On fail
	_ = oc.Producer.WriteMessages(context.Background(), kafka.Message{
		Topic: StorageCancelOrderTopic,
		Key:   []byte(order.Id.Hex()),
		Value: nil,
	})

	_, err := oc.UpdateOrderStatus(order.Id, models.OrderCancelPending)
	_, err = oc.UpdateOrderStatus(order.Id, models.OrderReservationCancelPending)
	if err != nil {
		oc.log.Error(err)
		return
	}

	// On success
	_, err = oc.UpdateOrderStatus(order.Id, models.OrderPaid)
	if err != nil {
		oc.log.Error(err)
		return
	}
}

func (oc *OrderCoordinator) OnOrderPayCanceled(order *models.Order) {
	order, err := oc.UpdateOrderStatus(order.Id, models.OrderReserved)
	if err != nil {
		oc.log.Error(err)
		return
	}

	_ = oc.Producer.WriteMessages(context.Background(), kafka.Message{
		Topic: StorageCancelOrderTopic,
		Key:   []byte(order.Id.Hex()),
		Value: nil,
	})
}*/
/*
func (oc *OrderCoordinator) HandleResponse(m *kafka.Message) {
	orderId, err := primitive.ObjectIDFromHex(string(m.Key))
	if err != nil {
		oc.log.Error(err)
		return
	}

	// TODO: move to repository
	order, err := oc.FindOrderId(orderId)
	if err != nil {
		oc.log.Error(err)
		return
	}

	// use mapping for resolving message destination

	switch m.Topic {
	case StorageReserveOrderResponseTopic:
		oc.OnOrderReserved(order)
	case StorageCancelOrderResponseTopic:
		oc.OnOrderReserveCanceled(order)
	case WalletPayOrderResponseTopic:
		oc.OnOrderPaid(order)
	case WalletCancelOrderResponseTopic:
		oc.OnOrderPayCanceled(order)
	default:
		oc.log.Error(`unknown topic`)
	}
}
*/ /*
// Commit set order status to 'Ready' and publish an event to kafka to commit changes in wallet and storage.
func (oc *OrderCoordinator) Commit() {

}

// FindOrderId searching for the order with provided id.
func (oc *OrderCoordinator) FindOrderId(id primitive.ObjectID) (*models.Order, error) {
	filter := bson.D{{"_id", id}}
	single := oc.Orders.FindOne(context.TODO(), filter)
	if err := single.Err(); err != nil {
		return nil, err
	}

	order := new(models.Order)
	if err := single.Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (oc *OrderCoordinator) UpdateOrderStatus(id primitive.ObjectID, status models.OrderStatus) (*models.Order, error) {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"status", status}}}}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	single := oc.Orders.FindOneAndUpdate(context.TODO(), filter, update, option)
	if err := single.Err(); err != nil {
		return nil, err
	}

	order := new(models.Order)
	if err := single.Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (oc *OrderCoordinator) UpdateOrder(id primitive.ObjectID, updates []bson.E) (*models.Order, error) {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", updates}}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	single := oc.Orders.FindOneAndUpdate(context.TODO(), filter, update, option)
	if err := single.Err(); err != nil {
		return nil, err
	}

	order := new(models.Order)
	if err := single.Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}
*/
