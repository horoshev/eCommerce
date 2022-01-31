package application

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
)

type WalletResources struct {
	Cfg           *Config
	Ctx           context.Context
	Log           *zap.SugaredLogger
	Mongo         *mongo.Client
	Database      *mongo.Database
	KafkaProducer *kafka.Writer
}

func NewWalletResources(ctx context.Context, log *zap.SugaredLogger, cfg *Config) *WalletResources {
	w := new(WalletResources)

	w.Cfg = cfg
	w.Ctx = ctx
	w.Log = log

	return w
}

func (w *WalletResources) Initialize() *WalletResources {

	uri := options.Client().ApplyURI(w.Cfg.MongoConnectionUrl)
	c, err := mongo.Connect(w.Ctx, uri)
	if err != nil {
		w.Log.Fatal(err)
	}

	err = c.Ping(context.Background(), nil)
	if err != nil {
		w.Log.Fatal(err)
	}

	w.Database = c.Database(ServiceName)
	w.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP(w.Cfg.KafkaConnectionUrl),
		Balancer: new(kafka.LeastBytes),
	}

	return w
}

func (w *WalletResources) Release(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := w.Database.Client().Disconnect(ctx); err != nil {
			w.Log.Error("Got an error while releasing mongodb resources.", "err", err)
		}
	}(ctx, &wg)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := w.KafkaProducer.Close(); err != nil {
			w.Log.Error("Got an error while releasing kafka resources.", "err", err)
		}
	}(ctx, &wg)

	wg.Wait()
}
