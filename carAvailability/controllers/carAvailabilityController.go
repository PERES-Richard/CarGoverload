package controllers

import (
	. "carAvailability/entities"
	"carAvailability/tools"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const MAX_SUPP_DURATION = 2

var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS"),
	Password: "",      // no password set
	DB:       redisDB, // use default DB
})

func readCarLockedByDay(yearDay int) []int {
	val, err := rdb.Get(context.Background(), strconv.Itoa(yearDay)).Result()
	log.Println("Checking car booked for day :", yearDay)
	log.Println("Value : ", val)
	if err == redis.Nil {
		log.Println("No cars booked for this date")
		return []int{}
	} else if err != nil {
		panic(err)
	}

	var carsLocked []int

	a := strings.Split(val, ",")
	for _, v := range a {
		var id int
		id, err = strconv.Atoi(v)
		carsLocked = append(carsLocked, id)
		if err != nil {
			log.Fatal("Error from redis: ", err)
		}
	}

	return carsLocked
}

func getNonAvailableCars(date time.Time) []int {
	carsAggregate := make([]int, 0)

	for i := 0; i < MAX_SUPP_DURATION; i++ {
		carsAggregate = append(readCarLockedByDay(date.YearDay()+i), carsAggregate...)
	}

	return carsAggregate
}

func checkIfCarsAreAvailable(cars []Car) bool {
	for _, car := range cars {
		log.Println("Car to check : ", car)
		log.Println("Day to check : ", car.DateDeparture.YearDay())
		bookedCarIds := readCarLockedByDay(car.DateDeparture.YearDay())
		for _, bookedCarId := range bookedCarIds {
			if bookedCarId == car.Id {
				return false
			}
		}
	}

	return true
}

func NewValidationSearchHandler(bookValidation BookValidationMessage, topic int) {
	result := checkIfCarsAreAvailable(bookValidation.Wishes)

	finalResult := BookValidationResult {
		WishId: bookValidation.WishId,
		IsValid: result,
	}

	resultJSON, err := json.Marshal(finalResult)
	if err != nil {
		log.Fatal("failed to marshal cars available:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), resultJSON) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

// Return the list of all car unavailable with given filters
func NewSearchHandler(message SearchMessage, topic int) {
	//date, err = time.Parse(time.RFC3339, dateParam[0])

	carsId := getNonAvailableCars(message.Date)
	log.Println("Booked cars : ", carsId)

	result := SearchResult{
		SearchId:     message.SearchId,
		CarsIdBooked: carsId,
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal("failed to marshal cars:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), resultJSON) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}
