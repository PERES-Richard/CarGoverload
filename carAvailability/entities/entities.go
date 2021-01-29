package entities

import "time"

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

// A Car representation for this svc
type Car struct {
	Id               int       `json:"id"`
	BeginBookedDate  time.Time `json:"beginBookedDate"`
	EndingBookedDate time.Time `json:"endingBookedDate"`
}

type SearchMessage struct {
	SearchId int       `json:"searchId"`
	Date     time.Time `json:"time"`
}

type SearchResult struct {
	SearchId     int   `json:"searchId"`
	CarsIdBooked []int `json:"carsIdBooked"`
}
