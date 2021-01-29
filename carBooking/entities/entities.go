package entities

import "time"

type CarBooked struct {
	Id               int       `json:"carId"`
	BeginBookedDate  time.Time `json:"beginBookedDate"`
	EndingBookedDate time.Time `json:"endingBookedDate"`
}
