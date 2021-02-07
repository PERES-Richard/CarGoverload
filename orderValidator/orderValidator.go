package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"orderValidator/tools"
	"os"
	"sync"

	controller "orderValidator/controllers"
	. "orderValidator/entities"
)

const VALIDATION_SEARCH_RESULT_READER_ID = 0
const BOOK_VALIDATION_TOPIC_READER_ID = 1
const BOOK_VALIDATION_RESULT_TOPIC_WRITER_ID = 0
const VALIDATION_SEARCH_WRITER_ID = 1

var readers = make([]*kafka.Reader, 2)
var wg sync.WaitGroup

func setUpKafka() {
	setupKafkaReaders()
	setupKafkaWriters()
}

func setupKafkaWriters() {
	configWriter := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "validation-search-requested",
		ClientId:  "orderValidator",
	}
	tools.SetUpWriter(VALIDATION_SEARCH_WRITER_ID, configWriter)

	configWriter = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-validation-result",
		ClientId:  "orderValidator",
	}
	tools.SetUpWriter(BOOK_VALIDATION_RESULT_TOPIC_WRITER_ID, configWriter)
}

func setupKafkaReaders() {
	// SetUpKafka
	configReader := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-validation",
		ClientId:  "orderValidator",
	}
	readers[BOOK_VALIDATION_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)

	configReader = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "validation-search-result",
		ClientId:  "orderValidator",
	}
	readers[VALIDATION_SEARCH_RESULT_READER_ID] = tools.GetUpKafkaReader(configReader)
}

func listenKafka(readerId int) {
	for {
		m, err := readers[readerId].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		messageHandlers(readerId, m)
	}
	wg.Done()
}

func messageHandlers(readerId int, m kafka.Message) {
	switch readerId {
	case BOOK_VALIDATION_TOPIC_READER_ID:
		{
			var parsedMessage BookMessage
			err := json.Unmarshal(m.Value, parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling book validation message:", err)
			}
			controller.BookValidationHandler(parsedMessage, VALIDATION_SEARCH_WRITER_ID)
		}
	case VALIDATION_SEARCH_RESULT_READER_ID:
		{
			var parsedMessage SearchResultMessage
			err := json.Unmarshal(m.Value, parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling search message:", err)
			}
			controller.ValidationSearchResultHandler(parsedMessage, BOOK_VALIDATION_RESULT_TOPIC_WRITER_ID)
		}
	}
}

func main() {
	// Setup readers & writers
	setUpKafka()

	wg.Add(2)

	go listenKafka(VALIDATION_SEARCH_RESULT_READER_ID)
	go listenKafka(BOOK_VALIDATION_TOPIC_READER_ID)

	wg.Wait()
}
