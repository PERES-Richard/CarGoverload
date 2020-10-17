package controllers

import (
	"bookingProcess/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func listSells(sellService *services.Service)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading sells"

		if err := json.NewEncoder(w).Encode(sellService.ListSells()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})

}

func createSell(sellService *services.Service)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error creating sell"
		var input struct {
			CustomerName    string `json:"customerName"`
			WagonType   string `json:"wagonType"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		tmp := sellService.CreateSell(input.CustomerName, input.WagonType, time.Now())

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(tmp); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})

}

func MakeSellHandlers(r *mux.Router, sellService *services.Service) {
	r.Handle("/booking-process/sells", listSells(sellService),
	).Methods("GET", "OPTIONS").Name("listSells")

	r.Handle("/booking-process/sells", createSell(sellService),
	).Methods("POST", "OPTIONS").Name("createSell")

}