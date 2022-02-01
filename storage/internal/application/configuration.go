package application

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName        string `envconfig:"SERVICE_NAME" default:"storage" required:"true"`
	MongoConnectionUrl string `envconfig:"MONGO_URL" default:"mongodb://root:password@mongodb-primary:27017/" required:"true"`
	KafkaConnectionUrl string `envconfig:"KAFKA_URL" default:"kafka:9092" required:"true"`
}

type ProductionConfig Config

type DevelopmentConfig struct {
	ServiceName        string `envconfig:"SERVICE_NAME" default:"storage" required:"true"`
	MongoConnectionUrl string `envconfig:"MONGO_URL" default:"mongodb://root:password@localhost:27017/" required:"true"`
	KafkaConnectionUrl string `envconfig:"KAFKA_URL" default:"localhost:9093" required:"true"`
}

func NewConfig() *Config {
	env := NewEnvironment()
	cfg, err := env.initConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

func (e Environment) initConfig() (cfg *Config, err error) {
	switch {
	case e.Is(Production):
		prod := ProductionConfig{}
		err = envconfig.Process("", &prod)
		cfg = (*Config)(&prod)
	default:
		dev := DevelopmentConfig{}
		err = envconfig.Process("", &dev)
		cfg = (*Config)(&dev)
	}

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
