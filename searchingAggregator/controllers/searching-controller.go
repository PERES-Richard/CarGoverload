package controllers

import (
	"context"
	"encoding/json"
	"log"
	. "searchingAggregator/entities"
	"searchingAggregator/tools"
)

const SEARCH_RESULT_TOPIC_WRITER_ID = 0
const VALIDATION_SEARCH_RESULT_TOPIC_WRITER_ID = 1

var searchMap = make(map[string]*SearchData)

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func AvailabilityResultHandler(parsedMessage AvailabilityResultMessage) {
	log.Println("Received availability results : ", parsedMessage.SearchId, "\nresults : ", parsedMessage.Cars)
	if searchMap[parsedMessage.SearchId] == nil {
		searchMap[parsedMessage.SearchId] = &SearchData{}
	}
	searchMap[parsedMessage.SearchId].AvailabilityResult = &parsedMessage.Cars
	checkResults(parsedMessage.SearchId)
}

func LocationResultHandler(parsedMessage LocationResultMessage) {
	log.Println("Received location results : ", parsedMessage.SearchId, "\nresults : ", parsedMessage.Cars)
	if searchMap[parsedMessage.SearchId] == nil {
		searchMap[parsedMessage.SearchId] = &SearchData{}
	}
	searchMap[parsedMessage.SearchId].LocationResult = &parsedMessage.Cars
	checkResults(parsedMessage.SearchId)
}

func NewSearchHandler(parsedMessage NewSearchMessage) {
	log.Println("Received new message : ", parsedMessage)
	if searchMap[parsedMessage.SearchId] == nil {
		searchMap[parsedMessage.SearchId] = &SearchData{}
	}
	searchMap[parsedMessage.SearchId].SearchTime = &parsedMessage.Date
	checkResults(parsedMessage.SearchId)
}

func NewValidationSearchHandler(parsedMessage NewSearchMessage) {
	log.Println("Received new message : ", parsedMessage)
	if searchMap[parsedMessage.SearchId] == nil {
		searchMap[parsedMessage.SearchId] = &SearchData{}
	}
	searchMap[parsedMessage.SearchId].SearchTime = &parsedMessage.Date
	searchMap[parsedMessage.SearchId].Validation = true
	checkResults(parsedMessage.SearchId)
}

func checkResults(searchId string) {
	log.Println("Checking results for searchId : ", searchId)
	if searchMap[searchId] != nil {
		if searchMap[searchId].SearchTime != nil && searchMap[searchId].LocationResult != nil && searchMap[searchId].AvailabilityResult != nil {
			endSearch(searchId)
		}
	}
}

func endSearch(searchId string) {
	log.Println("Ending search for searchId : ", searchId)

	var sd SearchData

	// if the parameter searchId does not correspond to an existing searchId we are waiting for
	if searchMap[searchId] == nil {
		return
	}

	sd = *searchMap[searchId]

	locationResults := sd.LocationResult
	availabilityResults := sd.AvailabilityResult

	// TODO: Enhance with multiple dates

	for _, car := range *locationResults {
		booked := false
		for _, bookedCar := range *availabilityResults {
			if bookedCar == car.Car.Id {
				booked = true
			}
		}
		if booked {
			log.Println("Car ", car.Car.Id, " is booked !")
			removeCar(locationResults, car.Car)
		}
	}

	if len(*locationResults) <= 0 {
		//TODO: relaunch request for compensation with another book date
	}

	offers := make([]Offer, 0)
	for _, car := range *locationResults {
		offers = append(offers, Offer{
			BookDate:  *sd.SearchTime,
			Arrival:   car.DestNode,
			Departure: car.Node,
			Car:       car.Car,
			Distance:  car.Distance,
		})
	}

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

	searchMap[searchId] = nil
}

func removeCar(carList *[]TrackedCar, car Car) {
	var result []TrackedCar
	for _, c := range *carList {
		if c.Car.Id != car.Id {
			result = append(result, c)
		}
	}
	*carList = result
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
