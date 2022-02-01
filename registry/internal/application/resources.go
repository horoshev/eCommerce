package application

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
)

type RegistryResources struct {
	cfg *Config
	ctx context.Context
	log *zap.SugaredLogger

	Producer *kafka.Writer
	Database *mongo.Database
}

func NewRegistryResources(ctx context.Context, log *zap.SugaredLogger, cfg *Config) *RegistryResources {
	resources := new(RegistryResources)
	resources.log = log
	resources.ctx = ctx
	resources.cfg = cfg

	return resources
}

func (r *RegistryResources) Initialize() *RegistryResources {
	r.Database = r.InitializeMongoDB()
	r.Producer = &kafka.Writer{
		Addr:     kafka.TCP(r.cfg.KafkaConnectionUrl),
		Balancer: new(kafka.LeastBytes),
	}

	return r
}

func (r *RegistryResources) InitializeMongoDB() *mongo.Database {
	uri := options.Client().ApplyURI(r.cfg.MongoConnectionUrl)
	c, err := mongo.Connect(r.ctx, uri)
	if err != nil {
		r.log.Fatal(err)
	}

	err = c.Ping(context.Background(), nil)
	if err != nil {
		r.log.Fatal(err)
	}

	return c.Database(ServiceName)
}

func (r *RegistryResources) Release(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := r.Database.Client().Disconnect(ctx); err != nil {
			r.log.Error("Got an error while releasing mongodb resources.", "err", err)
		}
	}(ctx, &wg)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := r.Producer.Close(); err != nil {
			r.log.Error("Got an error while releasing kafka resources.", "err", err)
		}
	}(ctx, &wg)

	wg.Wait()
}
