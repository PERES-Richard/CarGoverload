package main

import (
	"carAvailability/tools"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"

	. "carAvailability/entities"
)

const NEW_SEARCH_READER_ID = 0
const VALIDATION_SEARCH_READER_ID = 1
const CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID = 0

var readers = make([]*kafka.Reader, 2)

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS"),
	Password: "", // no password set
	DB:       0,  // use default DB
})

func readCarLocked(date time.Time) []Car {
	val, err := rdb.Get(context.Background(), string(date.YearDay())).Result()
	if err != nil {
		log.Panic("Error getting locked cars: ", err)
	}

	var carsLocked []Car
	err = json.Unmarshal([]byte(val), &carsLocked)
	if err != nil {
		log.Panic("Error unmarshaling search message:", err)
	}

	return carsLocked
}

// Filters & returns the list of all booked cars by filters
func getNonAvailableCars(date time.Time) []Car {
	var carsLocked []Car
	var carsLockedFiltered = make([]Car, 0)

	carsLocked = readCarLocked(date)

	for _, car := range carsLocked {
		if date.Before(car.BeginBookedDate) || date.After(car.EndingBookedDate) {
			carsLockedFiltered = append(carsLockedFiltered, car)
		}
	}

	return carsLockedFiltered
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
	tools.SetUpWriter(CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID,configWriter)
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

func NewValidationSearchHandler(message SearchMessage) {
	NewSearchHandler(message)
}

// Return the list of all car unavailable with given filters
func NewSearchHandler(message SearchMessage) {
	//date, err = time.Parse(time.RFC3339, dateParam[0])

	cars := getNonAvailableCars(message.Date)

	carsJSON, err := json.Marshal(cars)
	if err != nil {
		log.Fatal("failed to marshal cars:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID, []byte("value"), carsJSON) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:",kafkaErr)
	}
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
}

func messageHandlers(readerId int, m kafka.Message) {
	switch readerId {
	case NEW_SEARCH_READER_ID:
		{
			var parsedMessage SearchMessage
			err := json.Unmarshal(m.Value, parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling new search message:", err)
			}
			NewSearchHandler(parsedMessage)
		}
	case VALIDATION_SEARCH_READER_ID:
		{
			var parsedMessage SearchMessage
			err := json.Unmarshal(m.Value, parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling validation search message:", err)
			}
			NewValidationSearchHandler(parsedMessage)
		}
	}
}

func main() {
	// Setup readers & writers
	setUpKafka()

	go listenKafka(NEW_SEARCH_READER_ID)
	go listenKafka(VALIDATION_SEARCH_READER_ID)
}