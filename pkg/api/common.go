package api

import (
	"time"

	"github.com/philippheuer/go-kafka-consumer/config"
	"github.com/segmentio/kafka-go"
)

type MessageConsumer interface {
	ID() string
	Run() error
	Close()
}

func ReaderConfig(config *config.KafkaConfig, topic string, groupId string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:        []string{config.BootstrapServers},
		GroupID:        groupId,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	}
}
