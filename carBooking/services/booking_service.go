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

func (b *BookingService) CreateBook(date time.Time, dateArrival time.Time, car *entities.Car , supplier string, nodeDeparture *entities.Node, nodeArrival *entities.Node) entities.CarBooking {
	return repository.CreateBook(date, dateArrival, car, supplier, nodeDeparture, nodeArrival)
}

func (b *BookingService) FindAllBookings(typeId int) []entities.CarBooking{
	return repository.FindAllBookings(typeId)
}
