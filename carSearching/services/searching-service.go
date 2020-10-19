package services

import (
	"carSearching/entities"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type SearchingService struct {
	CAR_AVAILABILITY_URL string
}

func NewService() *SearchingService {
	var carURL string;
	if carURL = os.Getenv("CAR_AVAILABILITY_URL"); carURL == "" {
		carURL = "localhost:3001/car-availability"
		// OR raise error
	}
	return &SearchingService{
		CAR_AVAILABILITY_URL: carURL,
	}
}

func remove(s []entities.Car, i int) []entities.Car {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}


func (s *SearchingService) getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (s *SearchingService) Search(carType string, date string) []entities.Car{
	bookedCars, err := s.getBookedCars(carType, date)
	if err != nil {
		return []entities.Car{}
	}

	// carTracking service mocking
	// TODO: created carTracking service with mocking
	res := []entities.Car{entities.Car{Id: 1, CarType: entities.CarType{Name:"Liquid", Id:1}}, entities.Car{Id: 3, CarType: entities.CarType{Name:"Solid", Id:2}}}
	log.Println(res)

	return res //Todo Fix Length Error on the remove so returned res here
	// Remove booked cars from result
	for i, c := range res {
		b := false
		for _, bc := range bookedCars {
			if bc.Id == c.Id {
				b = true
			}
		}
		if b {
			res = remove(res, i)
		}
	}

	return res
}

// Get booked cars from carAvailability
func (s *SearchingService) getBookedCars(carType string, date string) ([]entities.Car, error) {
	var res []entities.Car
	err := s.getJson("http://" + s.CAR_AVAILABILITY_URL + "/getNonAvailableCars?carType=" + carType + "&date=" + date, &res)
	log.Println(res)
	return res, err
}
