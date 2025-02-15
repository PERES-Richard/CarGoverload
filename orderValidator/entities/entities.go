package entities

import "time"

type BookValidationResult struct {
	WishId		string		`json:"wishId"`
	IsValid		bool		`json:"isValid"`
}

type BookValidationMessage struct {
	Wishes		[]Wish	`json:"wishes"`
	WishId		string	`json:"wishId"`
}

type Wish struct {
	SearchId		string		`json:"searchId"`
	DepartureNode 	string    	`json:"departureNode"`
	ArrivalNode   	string    	`json:"arrivalNode"`
	DateDeparture 	time.Time 	`json:"dateDeparture"`
	DateArrival 	time.Time 	`json:"dateArrival"`
	CarId			int			`json:"carId"`
}

type BookConfirmation struct {
	WishId				string 			`json:"wishId"`
	Result				string			`json:"result"`
}
