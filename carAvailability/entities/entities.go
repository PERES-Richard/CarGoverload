package entities

import "time"

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

type BookValidationMessage struct {
	Wishes		[]Car	`json:"wishes"`
	WishId		string	`json:"wishId"`
}

type BookValidationResult struct {
	WishId		string	`json:"wishId"`
	IsValid		bool	`json:"isValid"`
}

// A Car representation for this svc
type Car struct {
	Id				int	    	`json:"carId"`
	DateDeparture 	time.Time	`json:"dateDeparture"`
	BookedYearDate	int
}

type SearchMessage struct {
	SearchId string    `json:"searchId"`
	Date     time.Time `json:"time"`
}

type SearchResult struct {
	SearchId     string		`json:"searchId"`
	CarsIdBooked []Car 		`json:"carsBookedByDay"`
}
