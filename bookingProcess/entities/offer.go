package entities

import "time"

type Offer struct {
	ID        		int      	`json:"id"`
	BookDate 		time.Time 	`json:"beginBookedDate"`
	Arrival 		Node		`json:"arrivalNode"`
	Departure 		Node		`json:"departureNode"`
	Car 			Car			`json:"car"`
	Price			float64 	`json:"price"`
	Duration 		int			`json:"duration"` //minutes
	BookArrival 	time.Time 	`json:"endingBookedDate"`
}
