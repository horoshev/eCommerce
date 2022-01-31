package consumers

import (
	"context"
	"eCommerce/wallet/internal/core"
	"eCommerce/wallet/internal/models"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"io"
	"log"
)

const (
	PayTopic = `wallet-pay-order`
	PayGroup = `wallet-pay-order-group`
)

// OrderConsumer for the user related events
type OrderConsumer struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	wallet *core.WalletController
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewOrderConsumer(ctx context.Context, log *zap.SugaredLogger, cfg *ConsumerConfig, wallet *core.WalletController) *OrderConsumer {
	consumer := new(OrderConsumer)

	consumer.ctx = ctx
	consumer.log = log
	consumer.reader = kafka.NewReader(cfg.ToReaderConfig(PayTopic, PayGroup))
	consumer.wallet = wallet

	return consumer
}

func (c *OrderConsumer) Start() {
	go func() {
		for {
			m, err := c.reader.ReadMessage(context.Background())
			if err == io.EOF {
				continue
			}

			if err != nil && err != io.EOF {
				c.log.Error(`read message err: `, err)
				continue
			}

			c.log.Info(`read message `, string(m.Key), m.Offset)
			_, err = c.ReserveCredit(m)
			if err != nil {
				continue
			}
		}
	}()
}

func (c *OrderConsumer) Stop() error {
	if err := c.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
		return err
	}

	return nil
}

// ReserveCredit event creates credit reservation for the customer.
func (c *OrderConsumer) ReserveCredit(message kafka.Message) (*models.Transaction, error) {
	order, err := ParseOrder(message)
	if err != nil {
		return nil, err
	}

	key, err := primitive.ObjectIDFromHex(string(message.Key))
	if err != nil {
		return nil, err
	}

	payment, err := c.wallet.PayOrder(key, order)
	if err != nil {
		return nil, err
	}

	// TODO: send response to registry

	return payment, nil
}

// CancelOrderTransaction commits reservation.
func (c *OrderConsumer) CancelOrderTransaction(message kafka.Message) error {
	order, err := ParseOrder(message)
	if err != nil {
		return err
	}

	refund, err := c.wallet.CancelOrder(order)
	if err != nil {
		return err
	}

	// TODO: ...
	_ = refund.Id
	return nil
}

func ParseOrder(message kafka.Message) (*models.Order, error) {
	order := new(models.Order)
	err := json.Unmarshal(message.Value, order)

	if err != nil {
		return nil, err
	}

	return order, nil
}
