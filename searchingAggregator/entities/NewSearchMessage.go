package entities

import "time"

type NewSearchMessage struct {
	SearchId 			int       	`json:"searchId"`
	Date     			time.Time 	`json:"time"`
	ExpectedResults 	int			`json:"expectedResults"`
}