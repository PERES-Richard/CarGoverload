package main

import (
	"carAvailability/tools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	. "carAvailability/entities"
)

const NEW_SEARCH_READER_ID = 0
const VALIDATION_SEARCH_READER_ID = 1
const CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID = 0

const MAX_SUPP_DURATION = 3

var readers = make([]*kafka.Reader, 2)

var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS"),
	Password: "",      // no password set
	DB:       redisDB, // use default DB
})

func readCarLockedByDay(yearDay int) []Car {
	val, err := rdb.Get(context.Background(), string(rune(yearDay))).Result()
	if err != nil {
		log.Panic("Error getting locked cars: ", err)
	}

	var carsLocked []Car

	a := strings.Split(val, ",")
	for _, v := range a {
		var id int
		id, err = strconv.Atoi(v)
		carsLocked = append(carsLocked, Car{
			Id:         id,
			BookedYearDate: yearDay,
		})
		if err != nil {
			log.Fatal("Error from redis: ", err)
		}
	}

	return carsLocked
}

func getNonAvailableCars(date time.Time) []Car {
	carsAggregate := make([]Car, 0)

	for i := 0; i < MAX_SUPP_DURATION; i++ {
		carsAggregate = append(readCarLockedByDay(date.YearDay()+i), carsAggregate...)
	}

	return carsAggregate
}

func NewValidationSearchHandler(message SearchMessage) {
	NewSearchHandler(message)
}

// Return the list of all car unavailable with given filters
func NewSearchHandler(message SearchMessage) {
	//date, err = time.Parse(time.RFC3339, dateParam[0])

	carsId := getNonAvailableCars(message.Date)

	result := SearchResult{
		SearchId:     message.SearchId,
		CarsIdBooked: carsId,
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal("failed to marshal cars:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), CAR_AVAILABILITY_RESULT_TOPIC_WRITER_ID, []byte("value"), resultJSON) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
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
