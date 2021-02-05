package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	. "searchingAggregator/entities"
	"searchingAggregator/tools"
	"time"
)

const SEARCH_RESULT_TOPIC_READER_ID = 0
const VALIDATION_SEARCH_RESULT_TOPIC_READER_ID = 1

const NO_SEARCH_ID = -1

var currentSearchId = NO_SEARCH_ID
var currentSearchDate time.Time
var ready = true
var locationResults []TrackedCar
var availabilityResults []int

var validationSearch = false

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func AvailabilityResultHandler(parsedMessage AvailabilityResultMessage) {
	if parsedMessage.SearchId == currentSearchId {
		availabilityResults = parsedMessage.Cars
		checkResults()
	} else {
		//TODO: produce compensation message
	}
}

func LocationResultHandler(parsedMessage LocationResultMessage) {
	if parsedMessage.SearchId == currentSearchId {
		locationResults = parsedMessage.Cars
		checkResults()
	} else {
		//TODO: produce compensation message
	}
}

func NewSearchHandler(parsedMessage NewSearchMessage) {
	if ready {
		ready = false
		currentSearchId = parsedMessage.SearchId
	} else {
		//TODO: produce compensation message
	}
}

func NewValidationSearchHandler(parsedMessage NewSearchMessage) {
	validationSearch = true
	NewSearchHandler(parsedMessage)
}

func checkResults() {
	if availabilityResults != nil && locationResults != nil {
		endSearch()
	}
}

func endSearch() {

	for _, car := range locationResults {
		booked := false
		for _, bookedCar := range availabilityResults {
			if bookedCar == car.Car.Id {
				booked = true
			}
		}
		if booked {
			locationResults, _ = removeCar(locationResults, car.Car)
		}
	}

	if len(locationResults) <= 0 {
		//TODO: relaunch request for compensation with another book date
	}

	offers := make([]Offer, 0)
	for _, car := range locationResults {
		offers = append(offers, Offer{
			BookDate:  currentSearchDate,
			Arrival:   car.DestNode,
			Departure: car.Node,
			Car:       car.Car,
		})
	}

	ready = true
	currentSearchId = NO_SEARCH_ID
	availabilityResults = nil
	locationResults = nil
	validationSearch = false

	resultJSON, err := json.Marshal(offers)
	if err != nil {
		log.Fatal("failed to marshal offers:", err)
		return
	}
	topic_id := SEARCH_RESULT_TOPIC_READER_ID
	if validationSearch {
		topic_id = VALIDATION_SEARCH_RESULT_TOPIC_READER_ID
	}
	kafkaErr := tools.KafkaPush(context.Background(), topic_id, []byte("value"), resultJSON)
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

func removeCar(carList []TrackedCar, car Car) ([]TrackedCar, error) {
	err := errors.New("Remove error: car not found")
	var result []TrackedCar
	for _, c := range carList {
		if c.Car.Id != car.Id {
			result = append(result, c)
		} else {
			err = errors.New("")
		}
	}
	return result, err
}

//func calculateDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
//	const PI float64 = 3.141592653589793
//
//	radlat1 := PI * lat1 / 180
//	radlat2 := PI * lat2 / 180
//
//	theta := lng1 - lng2
//	radtheta := PI * theta / 180
//
//	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta)
//
//	if dist > 1 {
//		dist = 1
//	}
//
//	dist = math.Acos(dist)
//	dist = dist * 180 / PI
//	dist = dist * 60 * 1.1515
//
//	if len(unit) > 0 {
//		if unit[0] == "K" {
//			dist = dist * 1.609344
//		} else if unit[0] == "N" {
//			dist = dist * 0.8684
//		}
//	}
//
//	return dist
//}
