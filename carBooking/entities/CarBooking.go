package entities

import "time"

type CarBooking struct {
	Supplier 		string		`json:"supplier"`
	Date  			time.Time		`json:"date"`
	Id 				int			`json:"id"`
	Arrival 		Node		`json:"arrivalNode"`
	Departure 		Node		`json:"departureNode"`
	Car 			Car			`json:"car"`
}
