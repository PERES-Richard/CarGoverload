package entities

import "time"

type CarBooking struct {
	Supplier 		string		`json:"supplier"`
	Date  			time.Time	`json:"date"`
	Id 				int			`json:"id" pg:",pk"`
	ArrivalId 		int			`json:"arrivalId"` //only used for pg orm
	Arrival 		*Node		`json:"arrivalNode" pg:"rel:has-one"`
	DepartureId 	int			`json:"departureId"` //only used for pg orm
	Departure 		*Node		`json:"departureNode" pg:"rel:has-one"`
	CarId		 	int			`json:"carId"` //only used for pg orm
	Car 			*Car		`json:"car" pg:"rel:has-one"`
}
