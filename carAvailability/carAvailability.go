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

	"github.com/robpike/filter"
)

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

// A Car representation for this svc
type Car struct {
	Id        int       `json:"id"`
	WagonType string    `json:"wagonType"` // TODO replace with enum ?
	Date      time.Time `json:"date"`
}

// Basic OK route for healthcheck
func ok(w http.ResponseWriter, req *http.Request) {
	_, err := io.WriteString(w, "ok")
	if err != nil {
		log.Fatal(err)
	}
}

// TODO Should return the list of car booked from DB and/or CarBooking service
func carsBookedList() []Car {
	carsBooked := make([]Car, 0)

	// TODO logic

	return carsBooked
}

// Filters & returns the list of all booked cars by filters
func getNonAvailableCars(date time.Time, wagonType string) []Car {
	var carsBookedFiltered []Car
	carsBooked := carsBookedList()

	// TODO logic
	var i interface{} = filter.Choose(carsBooked, func(car Car) bool {
		return car.WagonType == wagonType && car.Date.YearDay() == date.YearDay()
	})
	carsBookedFiltered, ok := i.([]Car)

	if !ok {
		log.Println("Error filtering booked cars")
	}

	return carsBookedFiltered
}

// Return the list of all car unavailable with given filters
func getNonAvailableCarsRoute(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()

	// Get the date from parameter
	dateParam, ok := params["date"]
	if !ok {
		log.Println("Error getNonAvailableCarsRoute : Date parameter not provided")
		return
	}
	// Convert DateParam into date
	date, err := time.Parse(time.RFC3339, dateParam[0])
	if err != nil {
		log.Println("Error getNonAvailableCarsRoute : Date parameter incorrect")
		log.Panic(err)
		return
	}

	// Get the wagonType from parameter
	wagonTypeParam, ok := params["wagonType"]
	if !ok {
		log.Fatalln("Error getNonAvailableCarsRoute : WagonType parameter not provided")
		return
	}
	wagonType := wagonTypeParam[0]

	cars := getNonAvailableCars(date, wagonType)

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

	// Create a new router to serve routes
	router := mux.NewRouter()

	// All the routes of the app
	router.HandleFunc("/car-availability/ok", ok).Methods("GET")
	router.HandleFunc("/car-availability/getNonAvailableCars", getNonAvailableCarsRoute).Methods("GET")

	fmt.Println("Server is running on port " + port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, router))
}
