package main

import (
	"carAvailability/tools"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"fmt"
	"log"
	"os"
	"time"

	. "carAvailability/entities"
)

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

// Return the list of all car unavailable with given filters
func NewSearchHandler(message SearchMessage) {
	//date, err = time.Parse(time.RFC3339, dateParam[0])

	cars := getNonAvailableCars(message.Date)

	carsJSON, err := json.Marshal(cars)
	if err != nil {
		log.Fatal("failed to marshal cars:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), []byte("value"), carsJSON) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:",kafkaErr)
	}
}

func main() {
	// SetUpKafka
	configReader := tools.KafkaConfig{
		BrokerUrl: "kafka-service:9092",
		Topic:     "new-search",
		ClientId:  "car-availability",
	}
	reader := tools.GetUpKafkaReader(configReader)
	defer reader.Close()

	configWriter := tools.KafkaConfig{
		BrokerUrl: "kafka-service:9092",
		Topic:     "car-availability-result",
		ClientId:  "car-availability",
	}
	tools.SetUpWriter(configWriter)

	// TODO Go func ?
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var parsedMessage SearchMessage
		err = json.Unmarshal(m.Value, parsedMessage)
		if err != nil {
			log.Panic("Error unmarshaling search message:", err)
		}

		NewSearchHandler(parsedMessage)
	}
}