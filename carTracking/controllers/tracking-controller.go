package controllers

import (
	"carTracking/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func GetCars(trackingService *services.TrackingService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := req.URL.Query()

		latParams, ok := params["latitude"]
		if !ok {
			log.Println("Error search : Node Id parameter not provided")
		}

		lgtParams, ok := params["longitude"]
		if !ok {
			log.Println("Error search : Node Id parameter not provided")
		}

		typeParams, ok := params["type"]
		if !ok {
			log.Println("Error search : Type parameter not provided")
		}

		res := trackingService.GetCars(latParams[0], lgtParams[0], typeParams[0])

		// Return search results as JSON Object
		jsonError := json.NewEncoder(w).Encode(res)
		if jsonError != nil {
			e := JSONError{Message: "Internal Server Error"}
			w.WriteHeader(http.StatusInternalServerError)
			err2 := json.NewEncoder(w).Encode(e)
			log.Panic(jsonError, err2)
		}
	})
}

func MakeTrackingHandlers(r *mux.Router, trackingService *services.TrackingService) {
	r.Handle("/car-tracking/get-cars", GetCars(trackingService)).Methods("GET", "OPTIONS").Name("getCars")
}