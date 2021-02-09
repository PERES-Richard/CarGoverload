// From https://yusufs.medium.com/getting-started-with-kafka-in-golang-14ccab5fa26

package tools

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	BrokerUrl     string
	Topic         string
	// ConsumerGroup string
	ClientId      string
}

func GetUpKafkaReader(kafkaConfig KafkaConfig) *kafka.Reader {
	brokers := strings.Split(kafkaConfig.BrokerUrl, ",")

	// make a new reader that consumes from topic-A
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaConfig.ClientId,
		Topic:           kafkaConfig.Topic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	return kafka.NewReader(config)
}