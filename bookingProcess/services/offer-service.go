package services

import (
	"bookingProcess/entities"
	"bookingProcess/utils"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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
		log.Println(err)
		return err
	}
	defer r.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	log.Println(buf.String())

	return json.Unmarshal([]byte(buf.String()), target)
}

func (s *OfferService) FindOffer(supplierName string, carType string, bookDate time.Time, departureNodeId string, arrivalNodeId string) ([]entities.Offer, error) {
	//Todo change results into DTO with car and nodes

	type SearchItem struct {
		BookDate 		time.Time 			`json:"bookDate"`
		Arrival 		entities.Node		`json:"arrivalNode"`
		Departure 		entities.Node		`json:"departureNode"`
		Car 			entities.Car		`json:"car"`
	}

	var results []SearchItem
	log.Println("Requeting sur carSearching")
	err := s.getJson("http://"+s.CAR_SEARCHING_HOST+":"+s.CAR_SEARCHING_PORT+"/car-searching/search?carType="+carType+"&date="+bookDate.Format(time.RFC3339)+"&departureNodeId="+departureNodeId+"&arrivalNodeId="+arrivalNodeId, &results)
	log.Println(results)

	var offers []entities.Offer

	for _, r := range results {
		kmDistance := s.determinePrice(r.Departure.Latitude, r.Departure.Longitude, r.Arrival.Latitude, r.Arrival.Longitude)

		number, err := strconv.Atoi(strconv.Itoa(rand.Int())[:8])
		if err != nil{
			log.Println(err)
			return []entities.Offer{}, err
		}

		duration := int(math.Floor(kmDistance*100)/100)

		offers = append(offers, entities.Offer{
			ID:        number,
			Arrival: r.Arrival,
			Departure: r.Departure,
			Car:    r.Car,
			BookDate:     r.BookDate,
			Price: math.Floor(kmDistance*2.5*100)/100,
			Duration: duration,
			BookArrival: r.BookDate.Add(time.Minute * time.Duration(duration)),
		})
	}

	found, supplier := s.findSupplierFromName(supplierName)
	if !found{
		supplierNew := entities.Supplier{
			ID:rand.Int(),
			Name:supplierName,
			Offers:[]entities.Offer{},
		}
		s.suppliers = append(s.suppliers, supplierNew)
		found, supplier = s.findSupplierFromName(supplierName)
	}

	supplier.Offers = append(supplier.Offers, offers...)
	log.Println("found supplier", supplier)
	log.Println("state of the object", s.suppliers)

	return offers, err
}

func (s *OfferService) findSupplierFromName(supplierName string) (bool, *entities.Supplier) {
	for i, n := range s.suppliers {
		if n.Name == supplierName {
			return true, &s.suppliers[i]

		}
	}

	return false, &entities.Supplier{}
}

func (s *OfferService) determinePrice(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := PI * lat1 / 180
	radlat2 := PI * lat2 / 180

	theta := lng1 - lng2
	radtheta := PI * theta / 180

	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	dist = dist * 1.609344


	return dist/ 3.3
}

func (s *OfferService) ListOffersOf(supplierName string) (error, []entities.Offer) {
	found, supplier := s.findSupplierFromName(supplierName)
	if found {
		return nil, supplier.Offers
	}
	return os.ErrNotExist, []entities.Offer{}
}

func (s *OfferService) PayOffer(id int, supplierName string) (error, entities.Offer) {
	found, supplier := s.findSupplierFromName(supplierName)
	if !found {
		return errors.New("Supplier not found : " + supplierName), entities.Offer{}
	}
	log.Println("Offer id passed : " +  strconv.Itoa(id))
	for _, offer := range supplier.Offers {
		log.Println("Offer id in supplier : " + strconv.Itoa(offer.ID))
		if offer.ID == id {
			done := s.bankAPI.PerformPayment(supplier.Name, offer.Price)
			if done{
				return nil, offer
			}
			return errors.New("Payment erreur"), entities.Offer{}
		}
	}

	return errors.New("Aucune recherche faite"), entities.Offer{}
}

func (s *OfferService) BookOffer(Ofr entities.Offer, supplierName string) interface{} {
	type SearchParams struct {
		Date string `json:"date"`
		CarId int `json:"carId"`
		Supplier string `json:"supplier"`
		NodeDepartureId int `json:"departureId"`
		NodeArrivalId int `json:"arrivalId"`
		DateArrival string `json:"dateArrival"`
	}

	var results struct {
		Supplier 		string				`json:"supplier"`
		Date  			time.Time			`json:"beginBookedDate"`
		DateArrival  	time.Time			`json:"endingBookedDate"`
		Id 				int					`json:"id"`
		Arrival 		entities.Node		`json:"arrivalNode"`
		Departure 		entities.Node		`json:"departureNode"`
		Car 			entities.Car		`json:"car"`
	}

	var body = SearchParams{
		Date:Ofr.BookDate.Format(time.RFC3339),
		CarId:Ofr.Car.Id,
		Supplier:supplierName,
		NodeArrivalId:Ofr.Arrival.Id,
		NodeDepartureId:Ofr.Departure.Id,
		DateArrival: Ofr.BookArrival.Format(time.RFC3339),
	}
	log.Println("Le body : ")
	log.Println(body)

	bodyByte, _ := json.Marshal(body)
	err := s.postJson("http://"+s.CAR_BOOKING_HOST+":"+s.CAR_BOOKING_PORT+"/car-booking/book", bytes.NewReader(bodyByte), &results)
	log.Println(results)
	if err != nil {
		log.Println(err)
		return err
	}
	return results
}
