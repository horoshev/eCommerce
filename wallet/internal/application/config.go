package application

import "eCommerce/wallet/internal/consumers"

type Environment string

const (
	Development Environment = `Development`
	Production  Environment = `Production`

	ServiceName = `wallet`
)

func (e Environment) Is(env Environment) bool {
	return e == env
}

type Config struct {
	ServiceName        string `envconfig:"SERVICE_NAME" default:"wallet" required:"true"`
	MongoConnectionUrl string `envconfig:"MONGO_URL" default:"mongodb://root:password@mongodb-primary:27017/" required:"true"`
	KafkaConnectionUrl string `envconfig:"KAFKA_URL" default:"kafka:9092" required:"true"`
}

type ProductionConfig Config

type DevelopmentConfig struct {
	ServiceName        string `envconfig:"SERVICE_NAME" default:"wallet" required:"true"`
	MongoConnectionUrl string `envconfig:"MONGO_URL" default:"mongodb://root:password@mongodb-primary:27017/" required:"true"`
	KafkaConnectionUrl string `envconfig:"KAFKA_URL" default:"kafka:9093" required:"true"`
}

func (c *Config) ToConsumerConfig() *consumers.ConsumerConfig {
	cfg := new(consumers.ConsumerConfig)
	cfg.KafkaConnectionUrl = c.KafkaConnectionUrl

	return cfg
}
