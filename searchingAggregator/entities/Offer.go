package entities

import "time"

type Offer struct {
	BookDate 		time.Time 	`json:"bookDate"`
	Arrival 		Node		`json:"arrivalNode"`
	Departure 		Node		`json:"departureNode"`
	Car 			Car			`json:"car"`
	Distance		float32		`json:"distance"`
}
