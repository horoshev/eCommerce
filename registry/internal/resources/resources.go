package resources

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type ConnectionConfig struct {
	MongoConnectionUrl string
	KafkaConnectionUrl string
}

type TopicConfig struct {
	OrdersTopic string
	UsersTopic  string
}

type Resources struct {
	//UsersProducer  *kafka.Writer
	//OrdersProducer *kafka.Writer
	Producer *kafka.Writer
	Mongo    *mongo.Database

	cfg *ConnectionConfig
	ctx context.Context
	log *zap.SugaredLogger
}

func New(ctx context.Context, log *zap.SugaredLogger, cfg *ConnectionConfig) *Resources {
	r := new(Resources)

	r.cfg = cfg
	r.ctx = ctx
	r.log = log

	return r
}

func (r *Resources) InitMongoDB() {
	uri := options.Client().ApplyURI(r.cfg.MongoConnectionUrl)
	client, err := mongo.Connect(r.ctx, uri)
	if err != nil {
		r.log.Fatal(err)
	}

	if err = client.Ping(r.ctx, nil); err != nil {
		r.log.Fatal(err)
	}

	r.Mongo = client.Database(`registry`)
}

func (r *Resources) InitKafka(cfg *TopicConfig) {
	kafkaURL := r.cfg.KafkaConnectionUrl
	r.Producer = &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Balancer: new(kafka.LeastBytes),
	}

	//
	//r.OrdersProducer = &kafka.Writer{
	//	Addr:     kafka.TCP(kafkaURL),
	//	Topic:    cfg.OrdersTopic,
	//	Balancer: new(kafka.LeastBytes),
	//}
	//
	//r.UsersProducer = &kafka.Writer{
	//	Addr:     kafka.TCP(kafkaURL),
	//	Topic:    cfg.UsersTopic,
	//	Balancer: new(kafka.LeastBytes),
	//}
}

func (r *Resources) Release() {
	if err := r.Mongo.Client().Disconnect(context.Background()); err != nil {
		r.log.Error(err)
	}

	if err := r.Producer.Close(); err != nil {
		r.log.Error(err)
	}
	//
	//if err := r.OrdersProducer.Close(); err != nil {
	//	r.log.Error(err)
	//}
	//
	//if err := r.UsersProducer.Close(); err != nil {
	//	r.log.Error(err)
	//}
}
