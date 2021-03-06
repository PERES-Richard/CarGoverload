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

var reader *kafka.Reader
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

//func arrayToString(a []int) string {
//	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]")
//}

func listenKafka() {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var parsedMessage WishBooked
		err = json.Unmarshal(m.Value, &parsedMessage)
		if err != nil {
			log.Panic("Error unmarshaling search message:", err)
		}

		BookRegisterHandler(parsedMessage)
	}
	wg.Done()
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

	wg.Add(1)
	go listenKafka()
	wg.Wait()
}
