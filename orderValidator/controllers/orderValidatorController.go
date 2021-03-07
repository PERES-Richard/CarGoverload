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

func ValidationSearchResultHandler(valid BookValidationResult, topicConfirmation int, topicRegister int) {
	if !valid.IsValid {
		toReturn := BookConfirmation {
			Result: "false",
			WishId: valid.WishId,
		}

		result, err := json.Marshal(toReturn)
		if err != nil {
			log.Fatal("failed to marshal cars:", err)
			return
		}

		kafkaErr := tools.KafkaPush(context.Background(), topicConfirmation, []byte("value"), result)
		if kafkaErr != nil {
			log.Panic("failed to write message:", kafkaErr)
		}
	} else {
		toReturn := BookValidationMessage {}

		for i := range bookingWaiting {
			if bookingWaiting[i].WishId == valid.WishId {
				toReturn = bookingWaiting[i]
			}
		}

		result, err := json.Marshal(toReturn)
		if err != nil {
			log.Fatal("failed to marshal cars:", err)
			return
		}

		kafkaErr := tools.KafkaPush(context.Background(), topicRegister, []byte("value"), result)
		if kafkaErr != nil {
			log.Panic("failed to write message:", kafkaErr)
		}
	}
}
