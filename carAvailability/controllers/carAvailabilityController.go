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

const MAX_SUPP_DURATION = 3

var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS"),
	Password: "",      // no password set
	DB:       redisDB, // use default DB
})

func readCarLockedByDay(yearDay int) []Car {
	val, err := rdb.Get(context.Background(), string(rune(yearDay))).Result()
	if err != nil {
		//log.Panic("Error getting locked cars: ", err)
		log.Println("No cars booked for this date")
		return []Car{}
	}

	var carsLocked []Car

	a := strings.Split(val, ",")
	for _, v := range a {
		var id int
		id, err = strconv.Atoi(v)
		carsLocked = append(carsLocked, Car{
			Id:             id,
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

func checkIfCarsAreAvailable(cars []Car) bool {
	for _, car := range cars {
		bookedCars := readCarLockedByDay(car.DateDeparture.YearDay())
		for _, bookedCar := range bookedCars {
			if bookedCar.Id == car.Id {
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
