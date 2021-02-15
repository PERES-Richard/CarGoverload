package entities

import (
	"encoding/json"
	"time"
)

type Node struct {
	Name      string      `json:"name"`
	Id        json.Number `json:"id,number"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
}

type Car struct {
	Id      int     `json:"id"`
	CarType CarType `json:"carType"`
}

type CarType struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Wish struct {
	DepartureNode string    `json:"departureNode"`
	ArrivalNode   string    `json:"arrivalNode"`
	DateDeparture time.Time `json:"dateDeparture"`
	CarType       string    `json:"carType"`
	NumberOfCars  int       `json:"numberOfCars"`
}

type Offer struct {
	BookDate  time.Time `json:"bookDate"`
	Arrival   Node      `json:"arrivalNode"`
	Departure Node      `json:"departureNode"`
	Car       Car       `json:"car"`
	Price     float32   `json:"price,omitempty"`
	//Score     float32   `json:"score,omitempty"`
}

type Book struct {
	Id         json.Number `json:"id,number"`
	Offers     []Offer     `json:"offers"`
	WishId     string      `json:"wishId"`
	TotalPrice float32     `json:"TotalPrice,omitempty"`
}
