// From https://yusufs.medium.com/getting-started-with-kafka-in-golang-14ccab5fa26

package tools

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var writer *kafka.Writer

func kafkaConfigure(kafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}

type KafkaConfig struct {
	BrokerUrl     string
	Topic         string
	// ConsumerGroup string
	ClientId      string
}

func SetUpKafka(kafkaConfig KafkaConfig) *kafka.Reader {
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

func KafkaPush(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		// TODO message
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return writer.WriteMessages(parent, message)
}

// TODO Kafka Pull
//for {
//m, err := reader.ReadMessage(context.Background())
//if err != nil {
//log.Error().Msgf("error while receiving message: %s", err.Error())
//continue
//}
//
//value := m.Value
//if m.CompressionCodec == snappy.NewCompressionCodec() {
//_, err = snappy.NewCompressionCodec().Decode(value, m.Value)
//}
//
//if err != nil {
//log.Error().Msgf("error while receiving message: %s", err.Error())
//continue
//}
//
//fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(value))
//}
//}