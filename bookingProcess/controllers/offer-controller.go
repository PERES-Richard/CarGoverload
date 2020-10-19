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

func listOffers(offerService *services.OfferService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading offers"

		vars := mux.Vars(r)
		tmp, _ := strconv.Atoi(vars["id"])

		if err := json.NewEncoder(w).Encode(offerService.ListOffersOf(tmp)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})

}

func findOffer(offerService *services.OfferService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error finding offer"
		var input struct {
			SupplierName    string `json:"supplierName"`
			CarType   string `json:"carType"`
			BookDate time.Time `json:"date"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		tmp, _ := offerService.FindOffer(input.SupplierName, input.CarType, input.BookDate)

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(tmp); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})

}

func payOffer(offerService *services.OfferService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Payment error"
		vars := mux.Vars(r)
		tmp, _ := strconv.Atoi(vars["id"])
		if err := json.NewEncoder(w).Encode(offerService.PayOffer(tmp)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})

}

func MakeOfferHandlers(r *mux.Router, offerService *services.OfferService) {
	r.Handle("/booking-process/suppliers/{id}/offers", listOffers(offerService),
	).Methods("GET", "OPTIONS").Name("listOffers")

	r.Handle("/booking-process/offers", findOffer(offerService),
	).Methods("POST", "OPTIONS").Name("findOffer")

	r.Handle("/booking-process/offers/{id}/payment", payOffer(offerService),
	).Methods("POST", "OPTIONS").Name("payOffer")

}