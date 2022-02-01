package application

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName        string `envconfig:"SERVICE_NAME" default:"wallet" required:"true"`
	MongoConnectionUrl string `envconfig:"MONGO_URL" default:"mongodb://root:password@mongodb-primary:27017/" required:"true"`
	KafkaConnectionUrl string `envconfig:"KAFKA_URL" default:"kafka:9092" required:"true"`
}

func NewConfig() *Config {
	cfg := new(Config)
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
