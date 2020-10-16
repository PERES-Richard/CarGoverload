package controllers

import (
	"bookingProcess/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func listSells()  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading sells"

		tmp := models.Sell{
			ID: 1,
			CustomerName: "UPS",
			Goods: "Liquide",
			BookDate: time.Now(),
		}
		if err := json.NewEncoder(w).Encode(tmp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})

}

func createSell()  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error creating sell"
		var input struct {
			CustomerName    string `json:"customerName"`
			Goods   string `json:"goodsType"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		tmp := models.Sell{
			ID:        0,
			CustomerName:     input.CustomerName,
			Goods:    input.Goods,
			BookDate:     time.Now(),
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(tmp); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})

}

func MakeSellHandlers(r *mux.Router) {
	r.Handle("/booking-process/sells", listSells(),
	).Methods("GET", "OPTIONS").Name("listSells")

	r.Handle("/booking-process/sells", createSell(),
	).Methods("POST", "OPTIONS").Name("createSell")

}