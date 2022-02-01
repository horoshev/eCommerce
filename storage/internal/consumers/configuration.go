package consumers

import (
	"github.com/segmentio/kafka-go"
)

var baseKafkaReaderConfiguration = kafka.ReaderConfig{
	MinBytes: 10e1,
	MaxBytes: 10e3,
}

func ReaderConfig(kafkaAddr, topic, group string) kafka.ReaderConfig {
	cfg := baseKafkaReaderConfiguration
	cfg.Brokers = []string{kafkaAddr}
	cfg.Topic = topic
	cfg.GroupID = group

	return cfg
}
