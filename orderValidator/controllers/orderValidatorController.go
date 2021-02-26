package controllers

import (
	"context"
	"encoding/json"
	"log"
	"orderValidator/tools"

	. "orderValidator/entities"
)

var bookingWaiting []BookValidationMessage

func BookValidationHandler(message BookValidationMessage, topic int) {
	bookingWaiting = append(bookingWaiting, message)

	wishData, err := json.Marshal(message)
	if err != nil {
		log.Fatal("failed to marshal wishes:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), wishData) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

func ValidationSearchResultHandler(valid BookValidationResult, topic int) {
	bookingWaitingFinal := make([]BookValidationMessage, 0)
	toReturn := BookValidationMessage {}

	for i := range bookingWaiting {
		if bookingWaiting[i].WishId == valid.WishId {
			toReturn = bookingWaiting[i]
		} else {
			bookingWaitingFinal = append(bookingWaitingFinal, toReturn)
		}
	}

	if !valid.IsValid {
		//TODO car is already booked by someone else so notify booking process api
	} else {
		result, err := json.Marshal(toReturn)
		if err != nil {
			log.Fatal("failed to marshal cars:", err)
			return
		}

		kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), result)
		if kafkaErr != nil {
			log.Panic("failed to write message:", kafkaErr)
		}
	}
}
