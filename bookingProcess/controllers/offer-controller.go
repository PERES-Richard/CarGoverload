package controllers

import (
	"bookingProcess/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func listOffers(offerService *services.OfferService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading offers"
		enableCors(&w)

		vars := mux.Vars(r)
		tmp, _ := vars["name"]

		error, offers := offerService.ListOffersOf(tmp)

		if error != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
		}

		if err := json.NewEncoder(w).Encode(offers); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})

}

func findOffer(offerService *services.OfferService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error finding offer"
		enableCors(&w)
		params := r.URL.Query()
		// Get the date from parameter
		supplier, ok := params["supplier"]
		if !ok {
			log.Println("supplier parameter not provided")
			return
		}
		// Get the date from parameter
		carTypeId, ok := params["carTypeId"]
		if !ok {
			log.Println("carTypeId parameter not provided")
			return
		}
		// Get the date from parameter
		arrivalNodeId, ok := params["arrivalNodeId"]
		if !ok {
			log.Println("arrivalNodeId parameter not provided")
			return
		}
		// Get the date from parameter
		departureNodeId, ok := params["departureNodeId"]
		if !ok {
			log.Println("departureNodeId parameter not provided")
			return
		}

		// Get the date from parameter
		dateTimeDeparture, ok := params["dateTimeDeparture"]
		if !ok {
			log.Println("dateTimeDeparture parameter not provided")
			return
		}
		// Convert DateParam into date
		date, err := time.Parse(time.RFC3339, dateTimeDeparture[0])
		if err != nil {
			log.Println("Date parameter incorrect")
			log.Panic(err)
			return
		}
		tmp, _ := offerService.FindOffer(supplier[0], carTypeId[0], date, departureNodeId[0], arrivalNodeId[0])

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
		enableCors(&w)
		var payParams struct{
			OfferId int `json:"offerId"`
			Supplier string `json:"supplier"`
		}
		_ = json.NewDecoder(r.Body).Decode(&payParams)
		payment, offer :=offerService.PayOffer(payParams.OfferId, payParams.Supplier)
		if !payment {
			log.Println(payment)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Payment error"))
			return
		}
		if err := json.NewEncoder(w).Encode(offerService.BookOffer(offer, payParams.Supplier)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Booking error"))
		}
	})

}

func MakeOfferHandlers(r *mux.Router, offerService *services.OfferService) {
	r.Handle("/booking-process/suppliers/{name}/offers", listOffers(offerService),
	).Methods("GET", "OPTIONS").Name("listOffers")

	r.Handle("/booking-process/offers", findOffer(offerService),
	).Methods("GET", "OPTIONS").Name("findOffer")

	r.Handle("/booking-process/offers/payment", payOffer(offerService),
	).Methods("POST", "OPTIONS").Name("payOffer")

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
