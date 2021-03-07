package main

import (
	. "carBooking/entities"
	"carBooking/tools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strconv"
	"sync"
)

var kafkaReader *kafka.Reader
var kafkaWriter *kafka.Writer
var wg sync.WaitGroup

var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS"),
	Password: "",      // no password set
	DB:       redisDB, // use default DB
})

var ctx = context.Background()

func BookRegisterHandler(wishBooked WishBooked) {
	log.Println("Storing in redis for book : ", wishBooked)
	bookByDay := make(map[string]string)
	for _, book := range wishBooked.CarsBooked {
		begin := book.BeginBookedDate.YearDay()
		end := book.EndingBookedDate.YearDay()

		log.Println("Begin : ", begin)
		log.Println("End : ", end)

		if bookByDay[strconv.Itoa(begin)] == "" {
			bookByDay[strconv.Itoa(begin)] =strconv.Itoa(book.CarId)
		} else {
			bookByDay[strconv.Itoa(begin)] = bookByDay[strconv.Itoa(begin)] + "," + strconv.Itoa(book.CarId)
		}

		if end != begin {
			if bookByDay[strconv.Itoa(end)] == "" {
				bookByDay[strconv.Itoa(end)] = strconv.Itoa(book.CarId)
			} else {
				bookByDay[strconv.Itoa(end)] = bookByDay[strconv.Itoa(end)] + "," + strconv.Itoa(book.CarId)
			}
		}
	}

	for key, value := range bookByDay {
		go lockCar(key, value)
	}

	bookConfirmation := BookConfirmation {
		WishId: wishBooked.WishId,
		Result: "true",
	}
	resultJSON, err := json.Marshal(bookConfirmation)
	if err != nil {
		log.Fatal("failed to marshal result:", err)
		return
	}

	kafkaErr := tools.KafkaPush(kafkaWriter, context.Background(), []byte("value"), resultJSON)
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

func lockCar(yearDay string, carIds string) {
	var cars string
	val, err := rdb.Get(ctx, yearDay).Result()
	if err == redis.Nil {
		log.Println("No car already booked")
		cars = carIds
		log.Println("Need to store : ", cars)
	} else if err != nil {
		panic(err)
	} else {
		log.Println("Some cars already booked")
		// Add the cars to the list
		cars = val + "," + carIds
		log.Println("Need to store : ", val, "and", carIds)
	}
	log.Println("Adding to database : ", cars, "for key", yearDay)
	err = rdb.Set(ctx, yearDay, cars, 0).Err()
	if err != nil {
		panic(err)
	}
}

func listenKafka() {
	for {
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var parsedMessage WishBooked
		err = json.Unmarshal(m.Value, &parsedMessage)
		if err != nil {
			log.Panic("Error unmarshaling search message:", err)
		}

		go BookRegisterHandler(parsedMessage)
	}
	wg.Done()
}

func setUpKafkaReader() {
	configReader := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-register",
		ClientId:  "car-booking",
	}
	kafkaReader = tools.GetUpKafkaReader(configReader)
}

func setupKafkaWriter() {
	configWriter := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-confirmation",
		ClientId:  "car-booking",
	}
	kafkaWriter = tools.GetKafkaWriter(configWriter)
}

func main() {
	// Setup readers & writers
	setUpKafkaReader()
	setupKafkaWriter()

	wg.Add(1)
	go listenKafka()
	wg.Wait()
}
