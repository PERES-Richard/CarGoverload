package main

import (
	"carSearching/entities"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var CAR_AVAILABILITY_URL string



// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func remove(s []entities.Car, i int) []entities.Car {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Get booked cars from carAvailability
func getBookedCars(carType string, date string) ([]entities.Car, error) {
	var res []entities.Car
	err := getJson("http://" + CAR_AVAILABILITY_URL + "/getNonAvailableCars?type=" + carType + "&date=" + date, res)
	return res, err
}

// Basic OK route for healthcheck
func ok(w http.ResponseWriter, req *http.Request) {
	_, err := io.WriteString(w, "ok")
	if err != nil {
		log.Fatal(err)
	}
}

func search(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := req.URL.Query()

	typeParams, ok := params["type"]
	if !ok {
		log.Println("Error search : Type parameter not provided")
	}

	dateParams, ok := params["date"]
	if !ok {
		log.Println("Error search : Date parameter not provided")
	}
	bookedCars, err := getBookedCars(typeParams[0], dateParams[0])
	if err != nil {
		e := JSONError{Message: "carAvailability service error"}
		w.WriteHeader(http.StatusInternalServerError)
		err2 := json.NewEncoder(w).Encode(e)
		log.Panic(err, err2)
	}

	// carTracking service mocking
	// TODO: created carTracking service with mocking
	res := []entities.Car{entities.Car{Id: 1, CarType: typeParams[0]}, entities.Car{Id: 3, CarType: typeParams[0]}}

	// Remove booked cars from result
	for i, c := range res {
		b := false
		for _, bc := range bookedCars {
			if bc.Id == c.Id {
				b = true
			}
		}
		if b {
			res = remove(res, i)
		}
	}

	// Return search results as JSON Object
	jsonError := json.NewEncoder(w).Encode(res)
	if jsonError != nil {
		e := JSONError{Message: "Internal Server Error"}
		w.WriteHeader(http.StatusInternalServerError)
		err2 := json.NewEncoder(w).Encode(e)
		log.Panic(jsonError, err2)
	}
}

func main() {
	// If there is a port variable set in env
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3003"
		// OR raise error
	}

	if CAR_AVAILABILITY_URL = os.Getenv("CAR_AVAILABILITY_URL"); CAR_AVAILABILITY_URL == "" {
		CAR_AVAILABILITY_URL = "localhost/car-availability"
		// OR raise error
	}

	// Create a new router to serve routes
	router := mux.NewRouter()

	// All the routes of the app
	router.HandleFunc("/car-searching/ok", ok).Methods("GET")

	router.HandleFunc("/car-searching/search", search).Methods("GET")

	fmt.Println("Server is running on port " + port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, router))
}
