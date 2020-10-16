package main

import (
	"carBooking/repository"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type JSONError struct {
	Message string `json:"Message"`
}

type SearchParams struct {
	Date string `json:"date"`
	CarId int `json:"carId"`
	Supplier string `json:"supplier"`
	NodeDepartureId int `json:"departureId"`
	NodeArrivalId int `json:"arrivalId"`
}

// Basic OK route for healthcheck
func ok(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "ok")
	if err != nil {
		log.Fatal(err)
	}
}

func findBookedCars(w http.ResponseWriter, _ *http.Request){
	w.Header().Set("Content-Type", "application/json")
	jsonError := json.NewEncoder(w).Encode(repository.FindAllBookings())
	if jsonError != nil {
		e := JSONError{Message: "Internal Server Error"}
		w.WriteHeader(http.StatusInternalServerError)
		err2 := json.NewEncoder(w).Encode(e)
		log.Panic(jsonError, err2)
	}
}

func bookCar(w http.ResponseWriter, req *http.Request){
	var sp SearchParams
	_ = json.NewDecoder(req.Body).Decode(&sp)

	nodeDeparture, err := repository.GetNodeFromId(sp.NodeDepartureId)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		_, err  := io.WriteString(w, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	nodeArrival, err := repository.GetNodeFromId(sp.NodeArrivalId)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		_, err  := io.WriteString(w, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	car, err := repository.GetCarFromId(sp.CarId)
	if err != nil{
		w.WriteHeader(http.StatusNotFound)
		_, err  := io.WriteString(w, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	date, err := time.Parse(time.RFC3339, sp.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err  := io.WriteString(w, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	//TODO check if car is not already used
	repository.CreateBook(date, car, sp.Supplier, nodeDeparture, nodeArrival)
	w.WriteHeader(http.StatusCreated)
	_, err = io.WriteString(w, "Ok it's booked")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// If there is a port variable set in env
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3002"
		// OR raise error
	}

	// Create a new router to serve routes
	router := mux.NewRouter()

	//TODO remove when bdd is up
	repository.InitMock()

	// All the routes of the app
	router.HandleFunc("/car-booking/ok", ok).Methods("GET")
	router.HandleFunc("/car-booking/findAll", findBookedCars).Methods("GET")
	router.HandleFunc("/car-booking/book", bookCar).Methods("POST")

	fmt.Println("Server is running on port " + port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, router))
}
