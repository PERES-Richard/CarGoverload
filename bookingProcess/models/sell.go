package models

import "time"

type Sell struct {
	ID        uint      `json:"id"`
	CustomerName      string    `json:"name"`
	Goods       string    `json:"goods"`
	BookDate time.Time `json:"bookDate"`
}