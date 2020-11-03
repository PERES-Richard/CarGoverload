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

func (b *BookingService) CreateBook(Date time.Time, Car *entities.Car , Supplier string, NodeDeparture *entities.Node, NodeArrival *entities.Node) entities.CarBooking {
	return repository.CreateBook(Date, Car, Supplier, NodeDeparture, NodeArrival)
}

func (b *BookingService) FindAllBookings(typeId int) []entities.CarBooking{
	return repository.FindAllBookings(typeId)
}

func (b *BookingService) GetAllTypes() []entities.CarType{
	return repository.GetAllTypes()
}
