package entities

import "time"

type Offer struct {
	ID        int      `json:"id"`
	BookDate time.Time `json:"bookDate"`
	Arrival 		int		`json:"arrivalNode"`
	Departure 		int		`json:"departureNode"`
	Car 			Car			`json:"car"`
	Price float32 `json:"price"`
}
