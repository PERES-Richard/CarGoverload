package controllers

import (
	"carSearching/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func search(searchingService *services.SearchingService)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := req.URL.Query()

		typeParams, ok := params["carType"]
		if !ok {
			log.Println("Error search : Type parameter not provided")
		}

		dateParams, ok := params["date"]
		if !ok {
			log.Println("Error search : Date parameter not provided")
		}
		date, err := time.Parse(time.RFC3339, dateParams[0])
		if err != nil {
			log.Println("Date parameter incorrect")
			log.Panic(err)
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

		res := searchingService.Search(typeParams[0], date, departureNodeId[0], arrivalNodeId[0])

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

func MakeSearchingHandlers(r *mux.Router, searchingService *services.SearchingService) {
	r.Handle("/car-searching/search", search(searchingService)).Methods("GET", "OPTIONS").Name("search")
}