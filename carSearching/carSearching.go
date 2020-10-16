package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var CAR_AVAILABILITY_URL string

type Car struct {
	Id int
	Type string
}

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func remove(s []Car, i int) []Car {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// Get booked cars from carAvailability
func getBookedCars(carType string, date string) ([]Car, error) {
	resp, err := http.Get(CAR_AVAILABILITY_URL + "/car-availability/getNonAvailableCars?type=" + carType + "&date=" + date)
	if err != nil {
		return []Car{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(body)
	var res []Car
	return res, nil /*err*/
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
	res := []Car{Car{Id: 1, Type: typeParams[0]}, Car{Id: 3, Type: typeParams[0]}}

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

	if port = os.Getenv("CAR_AVAILABILITY_URL"); port == "" {
		CAR_AVAILABILITY_URL = "localhost:3001"
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
