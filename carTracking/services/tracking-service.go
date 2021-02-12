package services

import (
	. "carTracking/entities"
	"strconv"
)

type TrackingService struct {
	Cars []Car
}

// MOCKING - IRL we would have an algorithm to get cars close to the location specified
// by the nodeId parameter and corresponding to the carType parameter,
// but for now we return a set list of cars
var cars = make([]Car, 100)

// Instantiate new service
func NewService() *TrackingService {
	// Mock the cars
	cars = append(cars, generateMockedCars(0, 5, "Solid", 1)...)
	cars = append(cars, generateMockedCars(50, 5, "Liquid", 2)...)

	return &TrackingService{
		Cars: cars,
	}
}

// Main search algorithm
func (s *TrackingService) GetCars(latitude string, longitude string, carType string) []Car {
	res := make([]Car, 0)
	for _, c := range cars {
		ct, _ := strconv.Atoi(carType)
		if c.CarType.Id == ct {
			res = append(res, c)
		}
	}
	return res
}

func generateMockedCars(startIndex int, numberOfCars int, typeName string, typeId int) (mocked []Car) {
	mocked = make([]Car, numberOfCars)
	for i := 0; i < numberOfCars; i++ {
		mocked = append(mocked, Car{
			Id: startIndex + i,
			CarType: CarType{
				Name: typeName,
				Id:   typeId,
			},
		})
	}
	return
}
