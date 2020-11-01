package entities

import "time"

// Custom error to return in case of a JSON parsing error
type JSONError struct {
Message string `json:"Message"`
}

// A Car representation for this svc
type Car struct {
Id      int     `json:"id"`
CarType CarType `json:"carType"`
Date    time.Time
}

// A Booking representation for this svc from carBooking
type Booking struct {
Supplier string    `json:"supplier"`
Date     time.Time `json:"date"`
Id       int       `json:"id"`
Car      Car       `json:"car"`
}

// A CarType representation for this svc from carBooking
type CarType struct {
Name string `json:"name"`
Id   int    `json:"id"`
}