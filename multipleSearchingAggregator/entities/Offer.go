package entities

import "time"

type Offer struct {
	DateDeparture time.Time `json:"dateDeparture"`
	Arrival       Node      `json:"arrivalNode"`
	Departure     Node      `json:"departureNode"`
	CarType       string    `json:"carType"`
	Car           Car    `json:"car,omitempty"`
	Distance      float32   `json:"distance"`
	NumberOfCars  int       `json:"numberOfCars"`
}
