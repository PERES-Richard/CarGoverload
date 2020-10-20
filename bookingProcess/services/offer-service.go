package services

import (
	"bookingProcess/entities"
	"bookingProcess/utils"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type OfferService struct {
	//repo Repository
	suppliers []entities.Supplier
	bankAPI utils.BankAPI

	CAR_SEARCHING_PORT string
	CAR_SEARCHING_HOST string

}

func NewService(suppliers []entities.Supplier) *OfferService {

	var carSrPort string;
	if carSrPort = os.Getenv("CAR_SEARCHING_PORT"); carSrPort == "" {
		carSrPort = "3003"
		// OR raise error
	}
	var carSrHost string;
	if carSrHost = os.Getenv("CAR_SEARCHING_HOST"); carSrHost == "" {
		carSrHost = "localhost"
		// OR raise error
	}
	return &OfferService{
		suppliers: suppliers,
		bankAPI: utils.BankAPI{
			Host: "localhost",
			Port:"9090",
			PaymentEP: "/pay",
		},
		CAR_SEARCHING_HOST:carSrHost,
		CAR_SEARCHING_PORT:carSrPort,

	}
}

func (s *OfferService) useBank(bank utils.BankAPI) {
	s.bankAPI = bank;
}

func (s *OfferService) getJson(url string, target interface{}) error {
	log.Println(url)
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (s *OfferService) FindOffer(supplierName string, carType string, bookDate time.Time) ([]entities.Offer, error) {

	var results []entities.Car
	err := s.getJson("http://"+s.CAR_SEARCHING_HOST+":"+s.CAR_SEARCHING_PORT+"/car-searching/search?carType="+carType+"&date="+bookDate.Format(time.RFC3339), &results)
	log.Println(results)

	var offers []entities.Offer

	for _, r := range results {
		offers = append(offers, entities.Offer{
			ID:        rand.Int(),
			Arrival: entities.Node{},
			Departure: entities.Node{},
			Car:    r,
			BookDate:     bookDate,
			Price: 0.0,
		})
	}

	log.Println("suppliernamevariable", supplierName)

	for _, n := range s.suppliers {
		if n.Name == supplierName {
			log.Println("found supplier")
			n.Offers = append(n.Offers, offers...)
			log.Println(n.Offers)


		}
	}
	return offers, err
}

func (s *OfferService) ListOffersOf(supplierId int) ([]entities.Offer) {
	for _, n := range s.suppliers {
		if n.ID == supplierId {
			return n.Offers

		}
	}
	return []entities.Offer{}
}

func (s *OfferService) PayOffer(id int) bool {

	for _, n := range s.suppliers {
		for _, i := range n.Offers {
			if i.ID == id {
				return s.bankAPI.PerformPayment(n.Name, i.Price)

			}
		}

	}

	return false
}