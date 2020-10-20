package services

import (
	"carBooking/entities"
	"carBooking/repository"
	"time"
)

type BookingService struct {

}

func NewService() *BookingService {
	//repository.InitMock()
	return &BookingService{

	}
}

func (b *BookingService) CreateBook(Date time.Time, Car entities.Car , Supplier string, NodeDeparture entities.Node, NodeArrival entities.Node) {
	repository.CreateBook(Date, Car, Supplier, NodeDeparture, NodeArrival)
}

func (b *BookingService) FindAllBookings() []entities.CarBooking{
	return repository.FindAllBookings()
}
