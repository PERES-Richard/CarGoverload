package services

import (
	"bookingProcess/models"
	"time"
)

type Service struct {
	//repo Repository
	currentSells []models.Sell
}

func NewService() *Service {
	return &Service{
		currentSells: make([]models.Sell,0),
	}
}

func (s *Service) CreateSell(customerName string, wagonType string, bookDate time.Time) (models.Sell) {
	sell := models.Sell{
		ID:        0,
		CustomerName:     customerName,
		WagonType:    wagonType,
		BookDate:     bookDate,
	}
	s.currentSells = append(s.currentSells, sell)
	return sell;
}

func (s *Service) ListSells() ([]models.Sell) {
	return s.currentSells;
}