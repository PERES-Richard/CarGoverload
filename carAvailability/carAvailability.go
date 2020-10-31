package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"robpike.io/filter"
)

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

// A Car representation for this svc
type Car struct {
	Id      int     `json:"id"`
	CarType CarType `json:"carType"`
	Date    time.Time
}

// A Booking representation for this svc from carBooking
type Booking struct {
	Supplier string    `json:"supplier"`
	Date     time.Time `json:"date"`
	Id       int       `json:"id"`
	Car      Car       `json:"car"`
}

// A CarType representation for this svc from carBooking
type CarType struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

var carBookingURL string
var GetBookingRoute string

// Basic OK route for healthcheck
func ok(w http.ResponseWriter, req *http.Request) {
	_, err := io.WriteString(w, "ok")
	if err != nil {
		log.Fatal(err)
	}
}

// Return the JSON data from the given URL
func getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Return the list of car booked from CarBooking service
func carsBookedListByType(carType string) []Car {
	bookings := make([]Booking, 0)
	err := getJson("http://"+carBookingURL+GetBookingRoute+carType, &bookings)
	if err != nil {
		log.Println(err)
	}
	log.Println(bookings)
	return bookingsToCars(bookings)
}

// Refactor Bookings list to Cars list
func bookingsToCars(bookings []Booking) []Car {
	cars := make([]Car, 0)

	for _, book := range bookings {
		cars = append(cars, Car{
			Id:      book.Car.Id,
			CarType: book.Car.CarType,
			Date:    book.Date,
		})
	}

	return cars
}

// Filters & returns the list of all booked cars by filters
func getNonAvailableCars(date time.Time, carType string) []Car {
	var carsBookedFiltered []Car
	carsBooked := carsBookedListByType(carType)

	var i interface{} = filter.Choose(carsBooked, func(car Car) bool {
		return car.Date.YearDay() == date.YearDay() /*&& car.CarType.Name == carType*/
	})
	carsBookedFiltered, ok := i.([]Car)

	if !ok {
		log.Println("Error filtering booked cars")
	}

	return carsBookedFiltered
}

// Return the list of all car unavailable with given filters
func GetNonAvailableCarsRoute(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()

	// Get the date from parameter
	dateParam, ok := params["date"]
	if !ok {
		log.Println("Error 1 GetNonAvailableCarsRoute : Date parameter not provided")
		return
	}
	// Convert DateParam into date
	date, err := time.Parse(time.RFC3339, dateParam[0])
	if err != nil {
		log.Println("Error 2 GetNonAvailableCarsRoute : Date parameter incorrect")
		log.Panic(err)
		return
	}

	// Get the carType from parameter
	carTypeParam, ok := params["carType"]
	if !ok {
		log.Fatalln("Error 3 GetNonAvailableCarsRoute : CarType parameter not provided")
		return
	}
	carType := carTypeParam[0]

	cars := getNonAvailableCars(date, carType)

	// Return logs as a JSON object
	jsonError := json.NewEncoder(w).Encode(cars)
	if jsonError != nil {
		e := JSONError{Message: "Internal Server Error"}
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(e)
		log.Panic(err)
	}

}

func main() {
	// If there is a port variable set in env
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3001"
		// OR raise error
	}

	var carBookingHost, carBookingPort string
	if carBookingHost = os.Getenv("CAR_BOOKING_HOST"); carBookingHost == "" {
		carBookingHost = "localhost"
	}

	if carBookingPort = os.Getenv("CAR_BOOKING_PORT"); carBookingPort == "" {
		carBookingPort = "3002"
	}

	carBookingURL = carBookingHost + ":" + carBookingPort

	if GetBookingRoute = os.Getenv("CARBOOKING_GETBOOKING_URL"); GetBookingRoute == "" {
		GetBookingRoute = "/car-booking/findAll/type/"
	}

	// Create a new router to serve routes
	router := mux.NewRouter()

	// All the routes of the app
	router.HandleFunc("/car-availability/ok", ok).Methods("GET")
	router.HandleFunc("/car-availability/getNonAvailableCars", GetNonAvailableCarsRoute).Methods("GET")

	fmt.Println("Server is running on port " + port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, router))
}
