package consumers

import (
	"context"
	"eCommerce/registry/internal/consumers/topics"
	"eCommerce/registry/internal/data"
	"eCommerce/registry/internal/models"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ConsumerBinding struct {
	topic   string
	handler func(message *kafka.Message)
}

func RunConsumers(baseConfig kafka.ReaderConfig, bindings []ConsumerBinding) {
	for _, bind := range bindings {
		baseConfig.Topic = bind.topic
		reader := kafka.NewReader(baseConfig)
		RunConsumer(reader, bind.handler)
	}
}

func RunConsumer(reader *kafka.Reader, handler func(message *kafka.Message)) {
	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				continue
			}

			handler(&m)
		}
	}()
}

type OrderConsumerSet struct {
	log        *zap.SugaredLogger
	repository data.RegistryRepository
	producer   *kafka.Writer

	bindings []ConsumerBinding
}

func NewOrderConsumerSet(log *zap.SugaredLogger, producer *kafka.Writer, repository data.RegistryRepository) *OrderConsumerSet {
	set := new(OrderConsumerSet)
	set.log = log
	set.producer = producer
	set.repository = repository

	set.bindings = []ConsumerBinding{
		{
			topic:   topics.StorageReserveOrderResponseTopic,
			handler: set.OrderReservedHandler,
		},
		{
			topic:   topics.StorageCancelOrderResponseTopic,
			handler: set.OrderReserveCanceledHandler,
		},
		{
			topic:   topics.WalletPayOrderResponseTopic,
			handler: set.OrderPaidHandler,
		},
		{
			topic:   topics.WalletCancelOrderResponseTopic,
			handler: set.OrderPayCanceledHandler,
		},
	}

	return set
}

func (oc *OrderConsumerSet) Run(baseCfg kafka.ReaderConfig) {
	RunConsumers(baseCfg, oc.bindings)
}

// OrderReservedHandler processing products reservation result.
// When products successfully reserved - update order status to 'ORDER_RESERVED', update a price of the order and
// finally publish an event to kafka.
// When products reservation failed - update order status to 'Error'.
func (oc *OrderConsumerSet) OrderReservedHandler(m *kafka.Message) {
	order, err := ParseOrder(m.Value)
	if err != nil {
		oc.log.Error(err)
		return
	}

	// TODO:
	// ~ 1. receive a price from storage
	// + 2. send price to reserve in wallet
	// + 3. add update order method. also can update a amount of order
	// + 4. add repository
	// order changes should be log

	if !IsSuccess(m) {
		order, err = oc.repository.UpdateOrderStatus(order.Id, models.OrderError)
		if err != nil {
			oc.log.Error(err)
		}
		return
	}

	order, err = oc.repository.UpdateOrder(order.Id, bson.D{
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

	err = oc.Publish(topics.WalletPayOrderTopic, order.Id.Hex(), value)
	if err != nil {
		oc.log.Error(err)
		return
	}

	order, err = oc.repository.UpdateOrderStatus(order.Id, models.OrderPaymentPending)
	if err != nil {
		oc.log.Error(err)
	}
}

func (oc *OrderConsumerSet) OrderReserveCanceledHandler(m *kafka.Message) {
	orderId, err := KeyOrderId(m)
	if err != nil {
		panic("")
	}

	if IsSuccess(m) {
		_, err = oc.repository.UpdateOrderStatus(orderId, models.OrderCanceled)
		if err != nil {
			oc.log.Error(err)
		}
		return
	}

	_, err = oc.repository.UpdateOrderStatus(orderId, models.OrderCancellationError)
	if err != nil {
		oc.log.Error(err)
	}
}

func (oc *OrderConsumerSet) OrderPaidHandler(m *kafka.Message) {
	orderId, err := KeyOrderId(m)
	if err != nil {
		oc.log.Error(err)
		return
	}

	if IsSuccess(m) {
		_, err := oc.repository.UpdateOrderStatus(orderId, models.OrderPaid)
		if err != nil {
			oc.log.Error(err)
		}
		return
	}

	_, err = oc.repository.UpdateOrderStatus(orderId, models.OrderCancelPending)
	if err != nil {
		oc.log.Error(err)
	}

	err = oc.Publish(topics.StorageCancelOrderTopic, orderId.Hex(), m.Value)
	if err != nil {
		oc.log.Error(err)
	}

	_, err = oc.repository.UpdateOrderStatus(orderId, models.OrderReservationCancelPending)
	if err != nil {
		oc.log.Error(err)
	}
}

func (oc *OrderConsumerSet) OrderPayCanceledHandler(m *kafka.Message) {
	orderId, err := KeyOrderId(m)
	if err != nil {
		oc.log.Error(err)
		return
	}

	if !IsSuccess(m) {
		_, err = oc.repository.UpdateOrderStatus(orderId, models.OrderCancellationError)
		if err != nil {
			oc.log.Error(err)
		}
		return
	}

	_, err = oc.repository.UpdateOrderStatus(orderId, models.OrderPaymentCanceled)
	if err != nil {
		oc.log.Error(err)
	}

	err = oc.Publish(topics.StorageCancelOrderTopic, orderId.Hex(), m.Value)
	if err != nil {
		oc.log.Error(err)
	}

	_, err = oc.repository.UpdateOrderStatus(orderId, models.OrderReservationCancelPending)
	if err != nil {
		oc.log.Error(err)
	}
}

func (oc *OrderConsumerSet) Publish(topic, key string, data []byte) error {
	return oc.producer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	})
}

func KeyOrderId(m *kafka.Message) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(string(m.Key))
}

func ParseOrder(data []byte) (*models.Order, error) {
	order := new(models.Order)
	err := json.Unmarshal(data, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func IsSuccess(m *kafka.Message) bool {
	for _, h := range m.Headers {
		if h.Key != `status` {
			continue
		}

		return len(h.Value) > 0 && h.Value[0] != 0
	}

	return false
}
