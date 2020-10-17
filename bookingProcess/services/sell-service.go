package services

import (
	"bookingProcess/models"
	"time"
)

type Service struct {
	//repo Repository
}

func NewService() *Service {
	return &Service{
	}
}

func (s *Service) CreateSell(customerName string, wagonType string, bookDate time.Time) (models.Sell) {
	sell := models.Sell{
		ID:        0,
		CustomerName:     customerName,
		WagonType:    wagonType,
		BookDate:     bookDate,
	}
	return sell;
}

func (s *Service) ListSells() ([]models.Sell) {
	var sells []models.Sell
	sells = append(sells, s.CreateSell("UPS", "Liquide", time.Now()))
	sells = append(sells, s.CreateSell("Amazon", "Solide", time.Now()))

	return sells;
}