package models

import "time"

type Sell struct {
	ID        uint      `json:"id"`
	CustomerName      string    `json:"customerName"`
	Goods       string    `json:"goods"`
	BookDate time.Time `json:"bookDate"`
}