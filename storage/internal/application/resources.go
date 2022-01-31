package application

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type StorageResources struct {
	appContext context.Context
	cfg        *Config
	log        *zap.SugaredLogger

	Database *mongo.Database
	Producer *kafka.Writer
}

func NewStorageResources(ctx context.Context, cfg *Config, log *zap.SugaredLogger) *StorageResources {
	sr := new(StorageResources)
	sr.appContext = ctx
	sr.cfg = cfg
	sr.log = log

	return sr
}

func (sr *StorageResources) Initialize() {
	sr.initDatabase()
	sr.initProducer()
}

func (sr *StorageResources) Release(ctx context.Context) error {
	if err := sr.Database.Client().Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (sr *StorageResources) initDatabase() {
	uri := options.Client().ApplyURI(sr.cfg.MongoConnectionUrl)
	c, err := mongo.Connect(sr.appContext, uri)

	if err != nil {
		sr.log.Fatal(err)
	}

	sr.Database = c.Database(sr.cfg.ServiceName)
}

func (sr *StorageResources) initProducer() {
	sr.Producer = &kafka.Writer{
		Addr:     kafka.TCP(sr.cfg.KafkaConnectionUrl),
		Balancer: new(kafka.LeastBytes),
	}
}
