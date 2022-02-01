package consumers

import (
	"context"
	"eCommerce/storage/internal/core"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"io"
	"log"
)

const (
	ReserveOrderTopic = `storage-reserve-order`
	ReserveOrderGroup = `storage-reserve-order-group`
	CancelOrderTopic  = `storage-cancel-order`
	CancelOrderGroup  = `storage-cancel-order-group`
)

type StorageConsumer struct {
	ctx     context.Context
	log     *zap.SugaredLogger
	storage core.StorageService

	reserveReader *kafka.Reader
	cancelReader  *kafka.Reader
}

func NewStorageConsumer(ctx context.Context, log *zap.SugaredLogger, kafkaAddr string, storage core.StorageService) *StorageConsumer {
	consumer := new(StorageConsumer)

	consumer.ctx = ctx
	consumer.log = log
	consumer.storage = storage

	consumer.reserveReader = kafka.NewReader(ReaderConfig(kafkaAddr, ReserveOrderTopic, ReserveOrderGroup))
	consumer.cancelReader = kafka.NewReader(ReaderConfig(kafkaAddr, CancelOrderTopic, CancelOrderGroup))

	return consumer
}

// ReserveOrder creates product reservation in storage.
func (c *StorageConsumer) ReserveOrder(message kafka.Message) error {
	order, err := ParseOrder(message)
	if err != nil {
		return err
	}

	err = c.storage.ReserveOrder(order)
	if err != nil {
		return err
	}

	return nil
}

// CancelOrder declines order products reservation.
func (c *StorageConsumer) CancelOrder(message kafka.Message) error {
	order, err := ParseOrder(message)
	if err != nil {
		return err
	}

	err = c.storage.CancelOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (c *StorageConsumer) Start() {
	c.launchConsumer(c.reserveReader, c.ReserveOrder)
	c.launchConsumer(c.cancelReader, c.CancelOrder)
}

func (c *StorageConsumer) Stop() error {
	if err := c.reserveReader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
		return err
	}

	if err := c.cancelReader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
		return err
	}

	return nil
}

func (c *StorageConsumer) launchConsumer(r *kafka.Reader, handler func(message kafka.Message) error) {
	go func(c StorageConsumer) {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				if err != io.EOF {
					c.log.Error(`read message err: `, err)
				}
				continue
			}

			c.log.Info(`read message`, string(m.Key), m.Offset)
			if err = handler(m); err != nil {
				c.log.Error(`handle message err:`, err)
				continue
			}
		}
	}(*c)
}
