package entities

import "time"

// TODO entities

type Wish struct {
	DepartureNode string    `json:"departureNode"`
	ArrivalNode   string    `json:"arrivalNode"`
	DateDeparture time.Time `json:"dateDeparture"`
	CarType       string    `json:"carType"`
	NumberOfCars  int       `json:"numberOfCars"`
}

type Offer struct {
	// TODO OFFER STRUCT
}
