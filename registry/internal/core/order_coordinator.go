package core

import (
	"context"
	"eCommerce/registry/internal/consumers"
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
		Topic: models.StorageReserveOrderTopic,
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
