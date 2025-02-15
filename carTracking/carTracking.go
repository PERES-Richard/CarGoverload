package main

import (
	"carTracking/controllers"
	"carTracking/services"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

// Basic OK route for healthcheck
func ok(w http.ResponseWriter, req *http.Request) {
	_, err := io.WriteString(w, "ok")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// If there is a port variable set in env
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3004"
		// OR raise error
	}

	// Create a new router to serve routes
	router := mux.NewRouter()

	trackingService := services.NewService()

	// All the routes of the app
	router.HandleFunc("/car-tracking/ok", ok).Methods("GET")
	controllers.MakeTrackingHandlers(router, trackingService)

	fmt.Println("Server is running on port " + port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, router))
}
