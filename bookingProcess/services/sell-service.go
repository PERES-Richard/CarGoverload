package services

import (
	"bookingProcess/entities"
	"bookingProcess/utils"
	"math/rand"
	"time"
)

type Service struct {
	//repo Repository
	currentSells []entities.Sell
	bankAPI utils.BankAPI
}

func NewService() *Service {
	return &Service{
		currentSells: make([]entities.Sell,0),
		bankAPI: utils.BankAPI{
			Host: "localhost",
			Port:"9090",
			PaymentEP: "/pay",
		},
	}
}

func (s *Service) useBank(bank utils.BankAPI) {
	s.bankAPI = bank;
}

func (s *Service) CreateSell(customerName string, wagonType string, bookDate time.Time, price float32) (entities.Sell) {
	sell := entities.Sell{
		ID:        rand.Int(),
		CustomerName:     customerName,
		WagonType:    wagonType,
		BookDate:     bookDate,
		Price: price,
	}
	s.currentSells = append(s.currentSells, sell)
	return sell;
}

func (s *Service) ListSells() ([]entities.Sell) {
	return s.currentSells;
}

func (s *Service) PaySell(id int) bool {

	for _, n := range s.currentSells {
		if n.ID == id {
			return s.bankAPI.PerformPayment(n.CustomerName, n.Price)

		}
	}

	return false
}