package entities

import "time"

type BookMessage struct {
	Data Wish
	OfferSelected Offer
}

type SearchResultMessage struct {
	OffersAvailable []Offer
}

type Wish struct {
	DepartureNode string    `json:"departureNode"`
	ArrivalNode   string    `json:"arrivalNode"`
	DateDeparture time.Time `json:"dateDeparture"`
	CarType       string    `json:"carType"`
	Amount        int       `json:"amount"`
}

type Offer struct {
	// TODO OFFER STRUCT
}