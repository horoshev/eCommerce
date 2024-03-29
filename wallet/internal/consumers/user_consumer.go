package consumers

import (
	"context"
	"eCommerce/wallet/internal/core"
	"eCommerce/wallet/internal/models"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
)

// UserConsumer for the user related events
type UserConsumer struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	wallet *core.WalletController
	reader *kafka.Reader
}

func NewUserConsumer(ctx context.Context, log *zap.SugaredLogger, kafkaAddr string, wallet *core.WalletController) *UserConsumer {
	consumer := new(UserConsumer)
	consumer.ctx = ctx
	consumer.log = log
	consumer.wallet = wallet

	readerConfig := ReaderConfig(kafkaAddr, models.WalletTopic, models.WalletGroup)
	consumer.reader = kafka.NewReader(readerConfig)

	return consumer
}

// NewWallet creates new wallet for the customer and initialize balance with some bonus.
func (u *UserConsumer) NewWallet(msg kafka.Message) (*models.Wallet, error) {
	user := new(models.User)
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return nil, err
	}

	wallet, err := u.wallet.CreateNewWallet(user)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (u *UserConsumer) Start() {
	go func() {
		for {
			m, err := u.reader.ReadMessage(u.ctx)
			if err != nil {
				continue
			}

			_, err = u.NewWallet(m)
			if err != nil {
				continue
			}
		}
	}()
}

func (u *UserConsumer) Stop() error {
	if err := u.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
		return err
	}

	return nil
}
