package entities

import "time"

type WishBooked struct {
	WishId           string        `json:"wishId"`
	CarsBooked		 []CarBook		`json:"wishes"`
}

type CarBook struct {
	SearchId			string			`json:"searchId"`
	DepartureNode		string			`json:"departureNode"`
	ArrivalNode			string			`json:"arrivalNode"`
	CarId				int				`json:"carId"`
	BeginBookedDate  	time.Time 		`json:"dateDeparture"`
	EndingBookedDate 	time.Time 		`json:"dateArrival"`
}
