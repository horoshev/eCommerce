package consumers

import (
	"github.com/segmentio/kafka-go"
	"log"
)

type ConsumerConfig struct {
	KafkaConnectionUrl string
}

func (cfg *ConsumerConfig) ToReaderConfig(topic, group string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:     []string{cfg.KafkaConnectionUrl},
		GroupID:     group,
		Topic:       topic,
		StartOffset: int64(0),
		MinBytes:    1,
		MaxBytes:    10e3,
		Logger:      log.Default(),
	}
}
