package model

type CarBooking struct {
	Supplier string
	Date  string
	Id int
	Arrival Node
	Departure Node
	Car Car
}
