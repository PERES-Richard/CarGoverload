package main

import (
	"carAvailability/tools"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"

	controller "carAvailability/controllers"
	. "carAvailability/entities"
)

const NEW_SEARCH_READER_ID = 0
const VALIDATION_SEARCH_READER_ID = 1
const CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID = 0
const VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID = 1

var readers = make([]*kafka.Reader, 2)
var wg sync.WaitGroup

func listenKafka(readerId int) {
	for {
		m, err := readers[readerId].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		go messageHandlers(readerId, m)
	}
	wg.Done()
}

func setUpKafka() {
	setupKafkaReaders()
	setupKafkaWriters()
}

func setupKafkaWriters() {
	configWriter := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "car-availability-result",
		ClientId:  "car-availability",
	}
	tools.SetUpWriter(CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID, configWriter)

	configWriter = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "validation-search-result",
		ClientId:  "car-availability",
	}
	tools.SetUpWriter(VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID, configWriter)
}

func setupKafkaReaders() {
	configReader := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "new-search",
		ClientId:  "car-availability",
	}
	readers[NEW_SEARCH_READER_ID] = tools.GetUpKafkaReader(configReader)

	configReader = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "validation-search",
		ClientId:  "car-availability",
	}
	readers[VALIDATION_SEARCH_READER_ID] = tools.GetUpKafkaReader(configReader)
}

func messageHandlers(readerId int, m kafka.Message) {
	switch readerId {
	case NEW_SEARCH_READER_ID:
		{
			var parsedMessage SearchMessage
			err := json.Unmarshal(m.Value, &parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling new search message:", err)
			}
			go controller.NewSearchHandler(parsedMessage, CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID)
		}
	case VALIDATION_SEARCH_READER_ID:
		{
			var parsedMessage BookValidationMessage
			err := json.Unmarshal(m.Value, &parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling validation search message:", err)
			}
			go controller.NewValidationSearchHandler(parsedMessage, VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID)
		}
	}
}

func main() {
	// Setup readers & writers
	setUpKafka()

	wg.Add(2)

	go listenKafka(NEW_SEARCH_READER_ID)
	go listenKafka(VALIDATION_SEARCH_READER_ID)

	wg.Wait()
}
