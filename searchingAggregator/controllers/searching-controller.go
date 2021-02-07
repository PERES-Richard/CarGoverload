package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	. "searchingAggregator/entities"
	"searchingAggregator/tools"
)

const SEARCH_RESULT_TOPIC_READER_ID = 0
const VALIDATION_SEARCH_RESULT_TOPIC_READER_ID = 1

var searchArrayList []SearchData

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func AvailabilityResultHandler(parsedMessage AvailabilityResultMessage) {
	if isSearchWaited(parsedMessage.SearchId) {
		for _, s := range searchArrayList {
			if s.SearchId == parsedMessage.SearchId {
				s.AvailabilityResult = parsedMessage.Cars
				checkResults(s.SearchId)
			}
		}
	} else {
		//TODO: produce compensation message
	}
}

func LocationResultHandler(parsedMessage LocationResultMessage) {
	if isSearchWaited(parsedMessage.SearchId) {
		for _, s := range searchArrayList {
			if s.SearchId == parsedMessage.SearchId {
				s.LocationResult = parsedMessage.Cars
				checkResults(s.SearchId)
			}
		}
	} else {
		//TODO: produce compensation message
	}
}

func NewSearchHandler(parsedMessage NewSearchMessage) {
	searchArrayList = append(searchArrayList, SearchData{
		SearchId: parsedMessage.SearchId,
		SearchTime: parsedMessage.Date,
		Validation: false,
	})
}

func NewValidationSearchHandler(parsedMessage NewSearchMessage) {
	searchArrayList = append(searchArrayList, SearchData{
		SearchId: parsedMessage.SearchId,
		SearchTime: parsedMessage.Date,
		Validation: true,
	})
}

func isSearchWaited(searchId string) bool {
	for _, s := range searchArrayList {
		if s.SearchId == searchId {
			return true
		}
	}
	return false
}

func checkResults(searchId string) {
	for _, s := range searchArrayList {
		if s.SearchId == searchId && s.AvailabilityResult != nil && s.LocationResult != nil {
			endSearch(searchId)
		}
	}
}

func endSearch(searchId string) {
	var sd SearchData
	found := false

	for _,s := range searchArrayList {
		if s.SearchId == searchId {
			sd = s
			found = true
		}
	}

	if !found {
		return
	}

	locationResults := sd.LocationResult
	availabilityResults := sd.AvailabilityResult

	// TODO: Enhance with multiple dates

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
			BookDate:  sd.SearchTime,
			Arrival:   car.DestNode,
			Departure: car.Node,
			Car:       car.Car,
		})
	}

	searchArrayList, _ = removeSearchData(sd.SearchId)

	resultJSON, err := json.Marshal(offers)
	if err != nil {
		log.Fatal("failed to marshal offers:", err)
		return
	}

	topic_id := SEARCH_RESULT_TOPIC_READER_ID
	if sd.Validation {
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

func removeSearchData(searchId string) ([]SearchData, error) {
	err := errors.New("Remove error: car not found")
	var result []SearchData
	for _, s := range searchArrayList {
		if s.SearchId != searchId {
			result = append(result, s)
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
