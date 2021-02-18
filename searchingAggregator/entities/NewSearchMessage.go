package entities

import "time"

type NewSearchMessage struct {
	SearchId 			string		`json:"searchId"`
	Date     			time.Time	`json:"dateDeparture"`
}