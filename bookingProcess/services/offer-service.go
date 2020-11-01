package services

import (
	"bookingProcess/entities"
	"bookingProcess/utils"
	"bytes"
	"encoding/json"
	"io"
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

	CAR_BOOKING_PORT string
	CAR_BOOKING_HOST string

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

	var carBkPort string;
	if carBkPort = os.Getenv("CAR_BOOKING_PORT"); carBkPort == "" {
		carBkPort = "3003"
		// OR raise error
	}
	var carBkHost string;
	if carBkHost = os.Getenv("CAR_BOOKING_HOST"); carBkHost == "" {
		carBkHost = "localhost"
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
		CAR_BOOKING_HOST:carBkHost,
		CAR_BOOKING_PORT:carBkPort,

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

func (s *OfferService) postJson(url string, body io.Reader,  target interface{}) error {
	log.Println(url)
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Post(url,"application/json", body )
	if err != nil {
		return err
	}
	defer r.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	return json.Unmarshal([]byte(buf.String()), target)
}

func (s *OfferService) FindOffer(supplierName string, carType string, bookDate time.Time, departureNodeId string, arrivalNodeId string) ([]entities.Offer, error) {

	var results []entities.Car
	log.Println("Requeting sur carSearching")
	err := s.getJson("http://"+s.CAR_SEARCHING_HOST+":"+s.CAR_SEARCHING_PORT+"/car-searching/search?carType="+carType+"&date="+bookDate.Format(time.RFC3339)+"&departureNodeId="+departureNodeId+"&arrivalNodeId="+arrivalNodeId, &results)
	log.Println(results)

	var offers []entities.Offer

	for _, r := range results {
		offers = append(offers, entities.Offer{
			ID:        rand.Int(),
			Arrival: entities.Node{Id:2},
			Departure: entities.Node{Id:1},
			Car:    r,
			BookDate:     bookDate,
			Price: 0.0,
		})
	}

	found, supplier := s.findSupplierFromName(supplierName)
	if !found{
		supplier = entities.Supplier{
			ID:rand.Int(),
			Name:supplierName,
			Offers:[]entities.Offer{},
		}
		s.suppliers = append(s.suppliers, supplier)
	}

	supplier.Offers = append(supplier.Offers, offers...)
	log.Println("found supplier", supplier)
	log.Println("state of the object", s.suppliers)

	return offers, err
}

func (s *OfferService) findSupplierFromName(supplierName string) (bool, entities.Supplier) {
	for i, n := range s.suppliers {
		if n.Name == supplierName {
			return true, s.suppliers[i]

		}
	}

	return false, entities.Supplier{}
}

func (s *OfferService) ListOffersOf(supplierName string) (error, []entities.Offer) {
	found, supplier := s.findSupplierFromName(supplierName)
	if found {
		return nil, supplier.Offers
	}
	return os.ErrNotExist, []entities.Offer{}
}

func (s *OfferService) PayOffer(id int, supplierName string) (bool, entities.Offer) {

	for _, n := range s.suppliers {
		for _, i := range n.Offers {
			if i.ID == id {
				return s.bankAPI.PerformPayment(n.Name, i.Price), i

			}
		}

	}

	return false, entities.Offer{}
}

func (s *OfferService) BookOffer(Ofr entities.Offer, supplierName string) interface{} {
	type SearchParams struct {
		Date string `json:"date"`
		CarId int `json:"carId"`
		Supplier string `json:"supplier"`
		NodeDepartureId int `json:"departureId"`
		NodeArrivalId int `json:"arrivalId"`
	}

	var results struct {
		Supplier 		string		`json:"supplier"`
		Date  			time.Time		`json:"date"`
		Id 				int			`json:"id"`
		Arrival 		entities.Node		`json:"arrivalNode"`
		Departure 		entities.Node		`json:"departureNode"`
		Car 			entities.Car			`json:"car"`
	}

	//Todo do real ID for car
	var body = SearchParams{Date:Ofr.BookDate.Format(time.RFC3339),CarId:1, Supplier:supplierName, NodeArrivalId:Ofr.Arrival.Id, NodeDepartureId:Ofr.Departure.Id}
	log.Println(body)

	bodyByte, _ := json.Marshal(body)
	err := s.postJson("http://"+s.CAR_BOOKING_HOST+":"+s.CAR_BOOKING_PORT+"/car-booking/book", bytes.NewReader(bodyByte), &results)
	log.Println(results)
	if err != nil {
		return err
	}
	return results
}