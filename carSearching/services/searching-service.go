package services

import (
	"carSearching/entities"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

type SearchingService struct {
	CAR_AVAILABILITY_PORT string
	CAR_AVAILABILITY_HOST string

	CAR_LOCATION_HOST string
	CAR_LOCATION_PORT string
}

// Instantiate new service
func NewService() *SearchingService {
	var carAvPort string;
	if carAvPort = os.Getenv("CAR_AVAILABILITY_PORT"); carAvPort == "" {
		carAvPort = "3001"
		// OR raise error
	}
	var carAvHost string;
	if carAvHost = os.Getenv("CAR_AVAILABILITY_HOST"); carAvHost == "" {
		carAvHost = "localhost"
		// OR raise error
	}

	var carLoPort string;
	if carLoPort = os.Getenv("CAR_LOCATION_PORT"); carLoPort == "" {
		carLoPort = "3005"
		// OR raise error
	}
	var carLoHost string;
	if carLoHost = os.Getenv("CAR_LOCATION_HOST"); carLoHost == "" {
		carLoHost = "localhost"
		// OR raise error
	}
	return &SearchingService{
		CAR_AVAILABILITY_PORT: carAvPort,
		CAR_AVAILABILITY_HOST: carAvHost,
		CAR_LOCATION_PORT:     carLoPort,
		CAR_LOCATION_HOST:     carLoHost,
	}
}

// Send request and store JSON result into target interface
func (s *SearchingService) sendRequest(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Main search algorithm
// 1: get unavailable cars at a certain date
// 2: Get nodes and cars close to the departure node
// 3: filter obtained cars with unavailable cars
// 4: create offers (associate cars to node and add prize)
func (s *SearchingService) Search(carType string, date time.Time, departureNodeId string, arrivalNodeId string) []entities.Offer{
	everyNodes, _ := s.getAllNodes()
	arrivalNode := s.getNodeFromId(everyNodes, arrivalNodeId)

	// Step 1: Get unavailable cars
	bookedCars, err := s.getBookedCars(carType, date)
	log.Println("Booked cars: ",bookedCars,"Err: ",err)
	if err != nil {
		return []entities.Offer{}
	}

	// Step 2 : Get every cars near given departure node and cartype
	trackedCars, err := s.getTrackedCarsAndNodes(carType, departureNodeId)
	log.Println("Tracked cars: ", trackedCars,"Err: ",err)
	if err != nil {
		return []entities.Offer{}
	}

	// Step 3 : Remove cars already booked
	for _, car := range trackedCars {
		booked := false
		for _, bookedCar := range bookedCars {
			if bookedCar.Id == car.Car.Id {
				booked = true
			}
		}
		if booked {
			trackedCars, _ = removeCar(trackedCars, car.Car)
		}
	}

	// Step 4 : return offers
	var offers = []entities.Offer{}
	for _, car := range trackedCars {
		price := s.calculateDistance(car.Node.Latitude, car.Node.Longitude, arrivalNode.Latitude, arrivalNode.Longitude, "K") / 3.3
		offers = append(offers, entities.Offer{
			BookDate:  date,
			Arrival:   arrivalNode,
			Departure: car.Node,
			Car:       car.Car,
			Price:     price,
		})
	}

	return offers
}

// Get booked cars from carAvailability
func (s *SearchingService) getBookedCars(carType string, date time.Time) ([]entities.Car, error) {
	res := make([]entities.Car, 0)
	query := "http://" + s.CAR_AVAILABILITY_HOST + ":" + s.CAR_AVAILABILITY_PORT + "/car-availability/getNonAvailableCars?carType=" + carType + "&date=" + date.Format(time.RFC3339)
	err := s.sendRequest(query, &res)
	log.Println(query + " ", res)
	return res, err
}

func (s *SearchingService) getAllNodes() ([]entities.Node, error) {
	res := make([]entities.Node, 0)
	err := s.sendRequest("http://" + s.CAR_LOCATION_HOST + ":" + s.CAR_LOCATION_PORT + "/car-location/findAllNodes", &res)
	log.Println(res)
	return res, err
}

func (s *SearchingService) getNodeFromId(nodes []entities.Node, id string) entities.Node {
	for _, node := range nodes {
		if node.Id == id{
			return node
		}
	}
	return entities.Node{}
}


// Get booked cars from carAvailability
func (s *SearchingService) getTrackedCarsAndNodes(carType string, nodeId string) ([]entities.TrackedCar, error) {
	res := make([]entities.TrackedCar, 0)
	err := s.sendRequest("http://" + s.CAR_LOCATION_HOST + ":" + s.CAR_LOCATION_PORT + "/car-location/searchTrackedCars?node=" + nodeId + "&carTypeId=" + carType + "&distance=100", &res)
	log.Println("Result of tracked cars and nodes", res)
	return res, err
}

// Remove element from array or splice
func removeCar(carList []entities.TrackedCar, car entities.Car) ([]entities.TrackedCar, error) {
	err := errors.New("Remove error: car not found")
	var result []entities.TrackedCar
	for _, c := range carList {
		if c.Car.Id != car.Id {
			result = append(result, c)
		} else {
			err = errors.New("")
		}
	}
	return result, err
}

func (s *SearchingService) calculateDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
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

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

