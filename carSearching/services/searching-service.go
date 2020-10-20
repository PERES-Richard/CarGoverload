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
	return &SearchingService{
		CAR_AVAILABILITY_PORT: carAvPort,
		CAR_AVAILABILITY_HOST: carAvHost,
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
func (s *SearchingService) Search(carType string, date string) []entities.Car{
	bookedCars, err := s.getBookedCars(carType, date)
	log.Println("Booked cars: ",bookedCars,"Err: ",err)
	if err != nil {
		return []entities.Car{}
	}

	// carTracking service mocking
	res := []entities.Car{entities.Car{Id: 1, CarType: entities.CarType{Name:"Liquid", Id:1}}, entities.Car{Id: 3, CarType: entities.CarType{Name:"Solid", Id:2}}}
	log.Println("Available cars, ",res)

	// Remove booked cars from result
	for _, car := range res {
		booked := false
		for _, bookedCar := range bookedCars {
			if bookedCar.Id == car.Id {
				booked = true
			}
		}
		if booked {
			res, _ = removeCar(res, car)
		}
	}
	return res
}

// Get booked cars from carAvailability
func (s *SearchingService) getBookedCars(carType string, date string) ([]entities.Car, error) {
	var res []entities.Car
	err := s.sendRequest("http://" + s.CAR_AVAILABILITY_HOST + ":" + s.CAR_AVAILABILITY_PORT + "/car-availability/getNonAvailableCars?carType=" + carType + "&date=" + date, &res)
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
