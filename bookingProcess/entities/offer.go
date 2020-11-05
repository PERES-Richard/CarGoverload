package entities

import "time"

type Offer struct {
	ID        int      `json:"id"`
	BookDate time.Time `json:"bookDate"`
	Arrival 		Node		`json:"arrivalNode"`
	Departure 		Node		`json:"departureNode"`
	Car 			Car			`json:"car"`
	Price float64 `json:"price"`

}
