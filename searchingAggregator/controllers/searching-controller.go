package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	. "searchingAggregator/entities"
	"searchingAggregator/tools"
)

const SEARCH_RESULT_TOPIC_WRITER_ID = 0
const VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID = 1

var searchArrayList []SearchData

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func AvailabilityResultHandler(parsedMessage AvailabilityResultMessage) {
	log.Println("Received availability results : ", parsedMessage.SearchId, "\nresults : ", parsedMessage.Cars)
	if !isSearchWaited(parsedMessage.SearchId) {
		searchArrayList = append(searchArrayList, SearchData{
			SearchId: parsedMessage.SearchId,
			Validation: false,
			ReceivedAvailability: false,
			ReceivedLocation: false,
			ReceivedTime: false,
		})
	}
	for i, _ := range searchArrayList {
		if searchArrayList[i].SearchId == parsedMessage.SearchId {
			searchArrayList[i].AvailabilityResult = parsedMessage.Cars
			searchArrayList[i].ReceivedAvailability = true
			checkResults(searchArrayList[i].SearchId)
		}
	}
}

func LocationResultHandler(parsedMessage LocationResultMessage) {
	log.Println("Received location results : ", parsedMessage.SearchId, "\nresults : ", parsedMessage.Cars)
	if !isSearchWaited(parsedMessage.SearchId) {
		searchArrayList = append(searchArrayList, SearchData{
			SearchId: parsedMessage.SearchId,
			Validation: false,
			ReceivedAvailability: false,
			ReceivedLocation: false,
			ReceivedTime: false,
		})
	}
	for i, _ := range searchArrayList {
		if searchArrayList[i].SearchId == parsedMessage.SearchId {
			searchArrayList[i].LocationResult = parsedMessage.Cars
			searchArrayList[i].ReceivedLocation = true
			checkResults(searchArrayList[i].SearchId)
		}
	}
}

func NewSearchHandler(parsedMessage NewSearchMessage) {
	log.Println("Received new message : ", parsedMessage)
	if !isSearchWaited(parsedMessage.SearchId) {
		searchArrayList = append(searchArrayList, SearchData{
			SearchId: parsedMessage.SearchId,
			SearchTime: parsedMessage.Date,
			Validation: false,
			ReceivedAvailability: false,
			ReceivedLocation: false,
			ReceivedTime: true,
		})
	} else {
		for i, _ := range searchArrayList {
			if searchArrayList[i].SearchId == parsedMessage.SearchId {
				searchArrayList[i].SearchTime = parsedMessage.Date
				searchArrayList[i].ReceivedTime = true
				checkResults(parsedMessage.SearchId)
			}
		}
	}
}

func NewValidationSearchHandler(parsedMessage NewSearchMessage) {
	if !isSearchWaited(parsedMessage.SearchId) {
		searchArrayList = append(searchArrayList, SearchData{
			SearchId: parsedMessage.SearchId,
			SearchTime: parsedMessage.Date,
			Validation: true,
			ReceivedAvailability: false,
			ReceivedLocation: false,
			ReceivedTime: true,
		})
	} else {
		for i, _ := range searchArrayList {
			if searchArrayList[i].SearchId == parsedMessage.SearchId {
				searchArrayList[i].SearchTime = parsedMessage.Date
				searchArrayList[i].ReceivedTime = true
				checkResults(parsedMessage.SearchId)
			}
		}
	}
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
	log.Println("Checking results for searchId : ", searchId)
	for _, s := range searchArrayList {
		log.Println(s, " - ", searchId)
		if s.SearchId == searchId && s.ReceivedAvailability && s.ReceivedLocation && s.ReceivedTime {
			endSearch(searchId)
		}
	}
}

func endSearch(searchId string) {
	log.Println("Ending search for searchId : ", searchId)

	var sd SearchData
	found := false

	for _,s := range searchArrayList {
		if s.SearchId == searchId {
			sd = s
			found = true
		}
	}

	// if the parameter searchId does not correspond to an existing searchId we are waiting for
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
			log.Println("Car ", car.Car.Id, " is booked !")
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

	result := ResultMessage{
		Offers:   offers,
		SearchId: searchId,
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal("failed to marshal result:", err)
		return
	}

	log.Println("Results for search ", searchId, " : ", offers)

	topic_id := SEARCH_RESULT_TOPIC_WRITER_ID
	if sd.Validation {
		topic_id = VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID
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
