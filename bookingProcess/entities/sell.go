package entities

import "time"

type Sell struct {
	ID        uint      `json:"id"`
	CustomerName      string    `json:"customerName"`
	WagonType       string    `json:"wagonType"`
	BookDate time.Time `json:"bookDate"`
	Price float32 `json:"price"`
}