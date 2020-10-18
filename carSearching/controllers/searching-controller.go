package controllers

import (
	"carSearching/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
		res := searchingService.Search(typeParams[0], dateParams[0])

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

	r.Handle("/car-searching/search", search(searchingService),).Methods("GET", "OPTIONS").Name("search")

}