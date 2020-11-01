package entities

import "time"

type CarBooking struct {
	Supplier 		string		`json:"supplier"`
	Date  			time.Time	`json:"date"`
	Id 				int64		`json:"id"`
	ArrivalId 		int64		`json:"arrivalId"` //only used for pg orm
	Arrival 		*Node		`json:"arrivalNode" pg:"rel:has-one"`
	DepartureId 	int64		`json:"departureId"` //only used for pg orm
	Departure 		*Node		`json:"departureNode" pg:"rel:has-one"`
	CarId		 	int64		`json:"carId"` //only used for pg orm
	Car 			*Car		`json:"car" pg:"rel:has-one"`
}
