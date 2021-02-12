package main

import (
	. "searchingAggregator/entities"
	"searchingAggregator/tools"
	"searchingAggregator/controllers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"sync"
)

const CAR_AVAILABILITY_RESULT_TOPIC_READER_ID = 0
const CAR_LOCATION_RESULT_TOPIC_READER_ID = 1
const NEW_SEARCH_TOPIC_READER_ID = 2
const VALIDATION_SEARCH_TOPIC_READER_ID = 3

const SEARCH_RESULT_TOPIC_WRITER_ID = 0
const VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID = 1

var readers = make([]*kafka.Reader, 4)
var wg sync.WaitGroup

func main() {
	// Setup readers & writers
	//setUpKafka()
	//
	//wg.Add(4)
	//
	//go listenKafka(CAR_AVAILABILITY_RESULT_TOPIC_READER_ID)
	//go listenKafka(CAR_LOCATION_RESULT_TOPIC_READER_ID)
	//go listenKafka(NEW_SEARCH_TOPIC_READER_ID)
	//go listenKafka(VALIDATION_SEARCH_TOPIC_READER_ID)
	//
	//wg.Wait()
	controllers.NewSearchHandler(NewSearchMessage{SearchId: "a"})
	controllers.LocationResultHandler(LocationResultMessage{
		SearchId: "a",
		Cars:     []TrackedCar{
			{
				Car:      Car{
					Id:      0,
					CarType: CarType{"liquid", 1},
				},
				Node:     Node{},
				DestNode: Node{},
			},
			{
				Car:      Car{
					Id:      1,
					CarType: CarType{"liquid", 1},
				},
				Node:     Node{},
				DestNode: Node{},
			},
			{
				Car:      Car{
					Id:     2,
					CarType: CarType{"liquid", 1},
				},
				Node:     Node{},
				DestNode: Node{},
			},
		},
	})
	controllers.AvailabilityResultHandler(AvailabilityResultMessage{SearchId: "a"})
}

func setUpKafka() {
	setupKafkaReaders()
	setupKafkaWriters()
}

func setupKafkaWriters() {
	configWriter := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "search-result",
		ClientId:  "searching-agregator",
	}
	tools.SetUpWriter(SEARCH_RESULT_TOPIC_WRITER_ID,configWriter)

	configWriter = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "validation-search-result",
		ClientId:  "searching-agregator",
	}
	tools.SetUpWriter(VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID,configWriter)
}

func setupKafkaReaders() {
	configReader := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "car-availability-result",
		ClientId:  "searching-agregator",
	}
	readers[CAR_AVAILABILITY_RESULT_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)

	configReader = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "car-location-result",
		ClientId:  "searching-agregator",
	}
	readers[CAR_LOCATION_RESULT_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)

	configReader = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "new-search",
		ClientId:  "searching-agregator",
	}
	readers[NEW_SEARCH_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)

	configReader = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "validation-search",
		ClientId:  "searching-agregator",
	}
	readers[VALIDATION_SEARCH_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)
}

func listenKafka(readerId int) {
	for {
		m, err := readers[readerId].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln("Error reader " + string(rune(readerId)) + ": ", err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		messageHandlers(readerId, m)
	}
	wg.Done()
}

func messageHandlers(readerId int, m kafka.Message) {
	switch readerId {
		case CAR_AVAILABILITY_RESULT_TOPIC_READER_ID:
			{
				var parsedMessage AvailabilityResultMessage
				err := json.Unmarshal(m.Value, &parsedMessage)
				if err != nil {
					log.Panic("Error unmarshalling availability result message:", err)
				}
				controllers.AvailabilityResultHandler(parsedMessage)
			}
		case CAR_LOCATION_RESULT_TOPIC_READER_ID:
			{
				var parsedMessage LocationResultMessage
				err := json.Unmarshal(m.Value, &parsedMessage)
				if err != nil {
					log.Panic("Error unmarshalling location result message:", err)
				}
				controllers.LocationResultHandler(parsedMessage)
			}
		case NEW_SEARCH_TOPIC_READER_ID:
			{
				var parsedMessage NewSearchMessage
				err := json.Unmarshal(m.Value, &parsedMessage)
				if err != nil {
					log.Panic("Error unmarshalling new search message:", err)
				}
				controllers.NewSearchHandler(parsedMessage)
			}
		case VALIDATION_SEARCH_TOPIC_READER_ID:
			{
				var parsedMessage NewSearchMessage
				err := json.Unmarshal(m.Value, &parsedMessage)
				if err != nil {
					log.Panic("Error unmarshalling validation search message:", err)
				}
				controllers.NewValidationSearchHandler(parsedMessage)
			}
	}
}