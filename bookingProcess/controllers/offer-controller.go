package controllers

import (
	"bookingProcess/entities"
	"bookingProcess/services"
	"bookingProcess/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"time"
)

const BOOK_VALIDATION_WRITER_ID = 0
const BOOK_REGISTER_WRITER_ID = 1
const WISH_REQUESTED_WRITER_ID = 2
const BOOK_VALIDATION_RESULT_READER_ID = 0
const WISH_RESULT_READER_ID = 1
var readers = make([]*kafka.Reader, 2)


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
		log.Println(payParams)
		err, offer :=offerService.PayOffer(payParams.OfferId, payParams.Supplier)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err := json.NewEncoder(w).Encode(offerService.BookOffer(offer, payParams.Supplier)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Booking error"))
		}
	})

}


func setUpKafka() {
	setupKafkaReaders()
	setupKafkaWriters()
}

func setupKafkaWriters() {
	configWriter := utils.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-validation",
		ClientId:  "booking-process",
	}
	utils.SetUpWriter(BOOK_VALIDATION_WRITER_ID,configWriter)

	configWriter = utils.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-register",
		ClientId:  "booking-process",
	}
	utils.SetUpWriter(BOOK_REGISTER_WRITER_ID,configWriter)

	configWriter = utils.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "wish-requested",
		ClientId:  "booking-process",
	}
	utils.SetUpWriter(WISH_REQUESTED_WRITER_ID,configWriter)
}

func setupKafkaReaders() {
	configReader := utils.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "book-validation-result",
		ClientId:  "booking-process",
	}
	readers[BOOK_VALIDATION_RESULT_READER_ID] = utils.GetUpKafkaReader(configReader)

	configReader = utils.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "validation-search",
		ClientId:  "booking-process",
	}
	readers[WISH_RESULT_READER_ID] = utils.GetUpKafkaReader(configReader)
}



func listenKafka(readerId int) {
	for {
		m, err := readers[readerId].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		messageHandlers(readerId, m)
	}
}

func messageHandlers(readerId int, m kafka.Message) {
	switch readerId {
	case BOOK_VALIDATION_RESULT_READER_ID:
		{
			var parsedMessage entities.SearchMessage
			err := json.Unmarshal(m.Value, parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling book validation message:", err)
			}
			//bookValidationHandler(parsedMessage) Todo do handler
		}
	case WISH_RESULT_READER_ID:
		{
			var parsedMessage entities.SearchMessage
			err := json.Unmarshal(m.Value, parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling validation search message:", err)
			}
			//wishResultHandler(parsedMessage) Todo do handler
		}
	}
}

func MakeOfferHandlers(r *mux.Router, offerService *services.OfferService) {
	r.Handle("/booking-process/suppliers/{name}/offers", listOffers(offerService),
	).Methods("GET", "OPTIONS").Name("listOffers")

	r.Handle("/booking-process/offers", findOffer(offerService),
	).Methods("GET", "OPTIONS").Name("findOffer")

	r.Handle("/booking-process/offers/payment", payOffer(offerService),
	).Methods("POST", "OPTIONS").Name("payOffer")

	// Setup readers & writers
	setUpKafka()

	go listenKafka(BOOK_VALIDATION_RESULT_READER_ID)
	go listenKafka(WISH_RESULT_READER_ID)


}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
