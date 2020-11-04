package services

import (
	"carSearching/entities"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type SearchingService struct {
	CAR_AVAILABILITY_PORT string
	CAR_AVAILABILITY_HOST string

	CAR_LOCATION_HOST string
	CAR_LOCATION_PORT string
}

// Instantiate new service
func NewService() *SearchingService {
	var carAvPort string;
	if carAvPort = os.Getenv("CAR_AVAILABILITY_PORT"); carAvPort == "" {
		carAvPort = "3001"
		// OR raise error
	}
	var carAvHost string;
	if carAvHost = os.Getenv("CAR_AVAILABILITY_HOST"); carAvHost == "" {
		carAvHost = "localhost"
		// OR raise error
	}

	var carLoPort string;
	if carLoPort = os.Getenv("CAR_LOCATION_PORT"); carLoPort == "" {
		carLoPort = "3005"
		// OR raise error
	}
	var carLoHost string;
	if carLoHost = os.Getenv("CAR_LOCATION_HOST"); carLoHost == "" {
		carLoHost = "localhost"
		// OR raise error
	}
	return &SearchingService{
		CAR_AVAILABILITY_PORT: carAvPort,
		CAR_AVAILABILITY_HOST: carAvHost,
		CAR_LOCATION_PORT:     carLoPort,
		CAR_LOCATION_HOST:     carLoHost,
	}
}

// Send request and store JSON result into target interface
func (s *SearchingService) sendRequest(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Main search algorithm
// 1: get unavailable cars at a certain date
// 2: Get nodes and cars close to the departure node
// 3: filter obtained cars with unavailable cars
// 4: create offers (associate cars to node and add prize)
func (s *SearchingService) Search(carType string, date time.Time, departureNodeId string, arrivalNodeId string) []entities.Car{
	// Step 1: Get unavailable cars
	bookedCars, err := s.getBookedCars(carType, date)
	log.Println("Booked cars: ",bookedCars,"Err: ",err)
	if err != nil {
		return []entities.Car{}
	}

	// TODO: Step 2 - get all necessary nodes

	// Step 3: carTracking service mocking
	// TODO: call for every necessary node
	//trackedCars, err := s.getTrackedCars(carType, departureNodeId)
	//log.Println("Tracked cars: ", trackedCars,"Err: ",err)
	//if err != nil {
	//	return []entities.Car{}
	//}
	//
	//// Step 4: Remove booked cars from result
	//for _, car := range trackedCars {
	//	booked := false
	//	for _, bookedCar := range bookedCars {
	//		if bookedCar.Id == car.Id {
	//			booked = true
	//		}
	//	}
	//	if booked {
	//		trackedCars, _ = removeCar(trackedCars, car)
	//	}
	//}
	//
	//// TODO - Step 5: create offers
	//
	//return trackedCars
	return nil
}

// Get booked cars from carAvailability
func (s *SearchingService) getBookedCars(carType string, date time.Time) ([]entities.Car, error) {
	res := make([]entities.Car, 0)
	query := "http://" + s.CAR_AVAILABILITY_HOST + ":" + s.CAR_AVAILABILITY_PORT + "/car-availability/getNonAvailableCars?carType=" + carType + "&date=" + date.Format(time.RFC3339)
	err := s.sendRequest(query, &res)
	log.Println(query + " ", res)
	return res, err
}

// Get booked cars from carAvailability
func (s *SearchingService) getTrackedCarsAndNodes(carType string, nodeId string) ([]entities.TrackedCar, error) {
	res := make([]entities.TrackedCar, 0)
	//TODO: remake request
	err := s.sendRequest("http://" + s.CAR_LOCATION_HOST + ":" + s.CAR_LOCATION_PORT + "/car-location/get-cars?nodeId=" + nodeId + "&type=" + carType, &res)
	log.Println(res)
	return res, err
}

// Remove element from array or splice
func removeCar(carList []entities.Car, car entities.Car) ([]entities.Car, error) {
	err := errors.New("Remove error: car not found")
	var result []entities.Car
	for _, c := range carList {
		if c.Id != car.Id {
			result = append(result, c)
		} else {
			err = errors.New("")
		}
	}
	return result, err
}
