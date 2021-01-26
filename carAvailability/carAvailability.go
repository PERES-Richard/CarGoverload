package main

import (
	"carAvailability/tools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	. "carAvailability/entities"
)

// If there is a port variable set in env
var port = os.Getenv("PORT")

// URL of the service
var CarBookingURL string

// URL to get all bookings
var GetBookingsRoute string

// URL to get bookings by type
var GetBookingsByTypeRoute string

// Return the JSON data from the given URL
func getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Return the list of car booked from CarBooking service
func bookingsByType(carType int) []Booking {
	bookings := make([]Booking, 0)

	getBookingsURL := "http://" + CarBookingURL

	// If there is no car type
	if carType == 0 {
		getBookingsURL += GetBookingsRoute
	} else {
		getBookingsURL += GetBookingsByTypeRoute + strconv.Itoa(carType)
	}

	// TODO should share the same BD with carBooking service
	err := getJson(getBookingsURL, &bookings)
	if err != nil {
		log.Println(err)
	}
	return bookings
}

// Filter bookings with filter func & refactor Bookings list to Cars list
func filterBookingsByFilter(bookings []Booking, filter func(car Car) bool) []Car {
	cars := make([]Car, 0)

	for _, book := range bookings {
		car := Car{
			Id:               book.Car.Id,
			CarTypeId:        book.Car.CarTypeId,
			BeginBookedDate:  book.BeginBookedDate,
			EndingBookedDate: book.EndingBookedDate,
		}

		if filter(car) {
			cars = append(cars, car)
		}
	}

	return cars
}

// Filters & returns the list of all booked cars by filters
func getNonAvailableCars(date time.Time, carType int) []Car {
	var carsBookedFiltered []Car
	bookings := bookingsByType(carType)
	log.Println("Les bookings au début")
	log.Println(bookings)

	var i interface{} = filterBookingsByFilter(bookings, func(car Car) bool {
		// If there is a date & the car is booked
		if date.Before(car.BeginBookedDate) || date.After(car.EndingBookedDate) {
			return false
		}

		// If there is a carType & the carType is different
		if car.CarTypeId != carType {
			return false
		}

		return true
	})
	carsBookedFiltered, ok := i.([]Car)
	log.Println(carsBookedFiltered)
	if !ok {
		log.Println("Error filtering booked cars")
	}

	log.Println("Les bookings après")
	log.Println(carsBookedFiltered)

	return carsBookedFiltered
}

// Return the list of all car unavailable with given filters
func GetNonAvailableCarsRoute(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()

	// Get the date from parameter
	dateParam, ok := params["date"]
	var date time.Time
	if ok {
		// Convert DateParam into date
		var err error
		date, err = time.Parse(time.RFC3339, dateParam[0])
		if err != nil {
			log.Println("Error 2 GetNonAvailableCarsRoute : Date parameter incorrect")
			log.Panic(err)
			return
		}
	}

	// Get the carType from parameter
	carTypeParam, _ := params["carType"]
	carType, _ := strconv.Atoi(carTypeParam[0])

	cars := getNonAvailableCars(date, carType)

	// TODO content
	kafkaErr := tools.KafkaPush(context.Background(), nil, []byte("your message content"))
	if kafkaErr != nil {
		log.Panic(kafkaErr)
	}

	// Return logs as a JSON object
	jsonError := json.NewEncoder(w).Encode(cars)
	if jsonError != nil {
		e := JSONError{Message: "Internal Server Error"}
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(e)
		log.Panic(err)
	}

}

func setupNetworking() {
	if port == "" {
		log.Fatal("No port set.")
	}

	var carBookingHost, carBookingPort string
	if carBookingHost = os.Getenv("CAR_BOOKING_HOST"); carBookingHost == "" {
		carBookingHost = "localhost"
	}

	if carBookingPort = os.Getenv("CAR_BOOKING_PORT"); carBookingPort == "" {
		carBookingPort = "3002"
	}

	CarBookingURL = carBookingHost + ":" + carBookingPort

	if GetBookingsByTypeRoute = os.Getenv("CARBOOKING_GETBOOKING_BY_TYPE_URL"); GetBookingsRoute == "" {
		GetBookingsByTypeRoute = "/car-booking/findAll/type/"
	}

	if GetBookingsRoute = os.Getenv("CARBOOKING_GETBOOKING_URL"); GetBookingsRoute == "" {
		GetBookingsRoute = "/car-booking/findAll/"
	}
}


func handleRoutes() (router *mux.Router){
	// Create a new router to serve routes
	router = mux.NewRouter()

	// All the routes of the app
	// Basic OK route for healthcheck
	router.HandleFunc("/car-availability/ok", func(w http.ResponseWriter, req *http.Request) { io.WriteString(w, "ok") }).Methods("GET")

	// Main handler
	router.HandleFunc("/car-availability/getNonAvailableCars", GetNonAvailableCarsRoute).Methods("GET")
	return
}

func main() {
	setupNetworking()

	// SetUpKafka
	config := tools.KafkaConfig{
		BrokerUrl: "kafka-service:9092",
		Topic:     "car-availability-result", // TODO voir multi topic
		ClientId:  "car-availability",
	}
	reader := tools.SetUpKafka(config)
	defer reader.Close()

	router := handleRoutes()

	fmt.Println("Server is running on port " + port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, router))
}