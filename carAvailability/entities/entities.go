package entities

import "time"

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

// A Car representation for this svc
type Car struct {
	Id               int `json:"id"`
	CarTypeId        int `json:"carType"`
	BeginBookedDate  time.Time
	EndingBookedDate time.Time
}

// A Booking representation for this svc from carBooking
type Booking struct {
	Supplier         string    `json:"supplier"`
	BeginBookedDate  time.Time `json:"beginBookedDate"`
	EndingBookedDate time.Time `json:"endingBookedDate"`
	Id               int       `json:"id"`
	Car              Car       `json:"carTypeId"`
}
