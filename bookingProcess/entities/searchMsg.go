package entities

import "time"

type SearchMessage struct {
	SearchId int       `json:"searchId"`
	Date     time.Time `json:"time"`
}
