package services

import (
	"carTracking/entities"
)

type TrackingService struct {
}

// Instantiate new service
func NewService() *TrackingService {
	return &TrackingService{
	}
}

// Main search algorithm
func (s *TrackingService) GetCars(latitude string, longitude string, carType string) []entities.Car{
	// MOCKING - IRL we would have an algorithm to get cars close to the location specified
	// by the nodeId parameter and corresponding to the carType parameter,
	// but for now we return a set list of cars
	return []entities.Car{
		{Id: 1, CarType: entities.CarType{Name: "Liquid", Id: 1}},
		{Id: 2, CarType: entities.CarType{Name: "Solid", Id: 2}},
		{Id: 3, CarType: entities.CarType{Name: "Solid", Id: 2}},
		{Id: 4, CarType: entities.CarType{Name: "Liquid", Id: 1}},
		{Id: 5, CarType: entities.CarType{Name: "Liquid", Id: 2}},
	}
}
