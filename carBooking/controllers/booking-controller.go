package controllers

import (
	"carBooking/repository"
	"carBooking/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
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
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		jsonError := json.NewEncoder(w).Encode(bookingService.FindAllBookings(-1))
		if jsonError != nil {
			e := JSONError{Message: "Internal Server Error"}
			w.WriteHeader(http.StatusInternalServerError)
			err2 := json.NewEncoder(w).Encode(e)
			log.Panic(jsonError, err2)
		}
	})

}

func findBookedCarsWithType(bookingService *services.BookingService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["typeId"], 10, 64)

		jsonError := json.NewEncoder(w).Encode(bookingService.FindAllBookings(id))
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
		enableCors(&w)
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

		jsonError := json.NewEncoder(w).Encode(bookingService.CreateBook(date, &car, sp.Supplier, &nodeDeparture, &nodeArrival))
		if jsonError != nil {
			e := JSONError{Message: "Internal Server Error"}
			w.WriteHeader(http.StatusInternalServerError)
			err2 := json.NewEncoder(w).Encode(e)
			log.Panic(jsonError, err2)
		}
	})
}

func getAllNodes(bookingService *services.BookingService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		jsonError := json.NewEncoder(w).Encode(bookingService.GetAllNodes())
		if jsonError != nil {
			e := JSONError{Message: "Internal Server Error"}
			w.WriteHeader(http.StatusInternalServerError)
			err2 := json.NewEncoder(w).Encode(e)
			log.Panic(jsonError, err2)
		}
	})
}

func getAllCarTypes(bookingService *services.BookingService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		jsonError := json.NewEncoder(w).Encode(bookingService.GetAllTypes())
		if jsonError != nil {
			e := JSONError{Message: "Internal Server Error"}
			w.WriteHeader(http.StatusInternalServerError)
			err2 := json.NewEncoder(w).Encode(e)
			log.Panic(jsonError, err2)
		}
	})
}

func MakeBookingHandlers(r *mux.Router, bookingService *services.BookingService) {
	r.Handle("/car-booking/getAllNodes", getAllNodes(bookingService),
	).Methods("GET", "OPTIONS").Name("getAllNodes")

	r.Handle("/car-booking/getAllCarTypes", getAllCarTypes(bookingService),
	).Methods("GET", "OPTIONS").Name("getAllCarTypes")

	r.Handle("/car-booking/findAll", findBookedCars(bookingService),
	).Methods("GET", "OPTIONS").Name("findAllBookings")

	r.Handle("/car-booking/findAll/type/{typeId}", findBookedCarsWithType(bookingService),
	).Methods("GET", "OPTIONS").Name("findAllBookingsForType")

	r.Handle("/car-booking/book", bookCar(bookingService),
	).Methods("POST", "OPTIONS").Name("bookCar")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
