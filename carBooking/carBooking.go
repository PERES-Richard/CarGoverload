package main

import (
	"carBooking/controllers"
	"carBooking/repository"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)



// Basic OK route for healthcheck
func ok(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "ok")
	if err != nil {
		log.Fatal(err)
	}
}



func main() {
	// If there is a port variable set in env
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3002"
		// OR raise error
	}

	// Create a new router to serve routes
	router := mux.NewRouter()

	//TODO remove when bdd is up
	repository.InitMock()

	// All the routes of the app
	router.HandleFunc("/car-booking/ok", ok).Methods("GET")
	controllers.MakeBookingHandlers(router)

	fmt.Println("Server is running on port " + port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, router))
}
