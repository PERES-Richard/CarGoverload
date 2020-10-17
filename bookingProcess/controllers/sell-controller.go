package controllers

import (
	"bookingProcess/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
		tmp := sellService.CreateSell(input.CustomerName, input.WagonType, time.Now(), 10.5)

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(tmp); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})

}

func paySell(sellService *services.Service)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading sells"
		vars := mux.Vars(r)
		tmp, _ := strconv.Atoi(vars["id"])
		if err := json.NewEncoder(w).Encode(sellService.PaySell(tmp)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})

}

func MakeSellHandlers(r *mux.Router, sellService *services.Service) {
	r.Handle("/booking-process/sells", listSells(sellService),
	).Methods("GET", "OPTIONS").Name("listSells")

	r.Handle("/booking-process/sells", createSell(sellService),
	).Methods("POST", "OPTIONS").Name("createSell")

	r.Handle("/booking-process/sell/{id}/payment", paySell(sellService),
	).Methods("POST", "OPTIONS").Name("paySell")

}