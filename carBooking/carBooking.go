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
	"strings"
)

var reader *kafka.Reader

var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS"),
	Password: "",      // no password set
	DB:       redisDB, // use default DB
})

var ctx = context.Background()

func BookRegisterHandler(carsBooked []CarBooked) {
	for _, book := range carsBooked {
		begin := book.BeginBookedDate.YearDay()
		end := book.EndingBookedDate.YearDay()

		for i := 0; i < end-begin; i++ {
			go lockCar(begin+i, book.Id)
		}
	}
}

func lockCar(yearDay int, carId int) {
	var cars, adding string
	val, err := rdb.Get(ctx, string(yearDay)).Result()
	if err == redis.Nil {
		// No cars previously booked
		adding = string(rune(carId))
	} else if err != nil {
		panic(err)
	} else {
		// Add the cars to the list
		cars = val
		adding = "," + string(rune(carId))
	}
	err = rdb.Set(ctx, string(yearDay), cars+adding, 0).Err()
	if err != nil {
		panic(err)
	}
}

func arrayToString(a []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]")
}

func listenKafka() {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var parsedMessage []CarBooked
		err = json.Unmarshal(m.Value, parsedMessage)
		if err != nil {
			log.Panic("Error unmarshaling search message:", err)
		}

		BookRegisterHandler(parsedMessage)
	}
}

func setUpKafka() {
	configReader := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-register",
		ClientId:  "car-booking",
	}
	reader = tools.GetUpKafkaReader(configReader)
}

func main() {
	// Setup readers & writers
	setUpKafka()

	go listenKafka()
}
