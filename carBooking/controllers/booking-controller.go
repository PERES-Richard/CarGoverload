package controllers

import (
	"carBooking/repository"
	"carBooking/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
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

func findBookedCars(bookingService *services.BookingService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonError := json.NewEncoder(w).Encode(bookingService.FindAllBookings())
		if jsonError != nil {
			e := JSONError{Message: "Internal Server Error"}
			w.WriteHeader(http.StatusInternalServerError)
			err2 := json.NewEncoder(w).Encode(e)
			log.Panic(jsonError, err2)
		}
	})

}

func bookCar(bookingService *services.BookingService)  http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sp SearchParams
		_ = json.NewDecoder(r.Body).Decode(&sp)

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
		bookingService.CreateBook(date, car, sp.Supplier, nodeDeparture, nodeArrival)
		w.WriteHeader(http.StatusCreated)
		_, err = io.WriteString(w, "Ok it's booked")
		if err != nil {
			log.Fatal(err)
		}}	)

}


func MakeBookingHandlers(r *mux.Router, bookingService *services.BookingService) {
	r.Handle("/car-booking/findAll", findBookedCars(bookingService),
	).Methods("GET", "OPTIONS").Name("findAllBookings")

	r.Handle("/car-booking/book", bookCar(bookingService),
	).Methods("POST", "OPTIONS").Name("bookCar")


}
