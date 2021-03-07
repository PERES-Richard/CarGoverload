// From https://yusufs.medium.com/getting-started-with-kafka-in-golang-14ccab5fa26

package tools

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var writer = make([]*kafka.Writer, 3)

func SetUpWriter(readerId int, kafkaConfig KafkaConfig) {
	brokers := strings.Split(kafkaConfig.BrokerUrl, ",")

	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: kafkaConfig.ClientId,
	}

	config := kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            kafkaConfig.Topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	writer[readerId] = kafka.NewWriter(config)
}

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

func KafkaPush(parent context.Context, writerId int, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return writer[writerId].WriteMessages(parent, message)
}
