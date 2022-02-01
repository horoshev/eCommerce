package consumers

import (
	"github.com/segmentio/kafka-go"
)

var baseConfiguration = kafka.ReaderConfig{
	MinBytes: 10e1,
	MaxBytes: 10e3,
}

func ReaderConfig(addr, topic, group string) kafka.ReaderConfig {
	cfg := baseConfiguration
	cfg.Brokers = []string{addr}
	cfg.Topic = topic
	cfg.GroupID = group

	return cfg
}
