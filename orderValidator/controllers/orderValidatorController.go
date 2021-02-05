package controllers

import (
	"context"
	"encoding/json"
	"log"
	"orderValidator/tools"

	. "orderValidator/entities"
)

func BookValidationHandler(message BookMessage, topic int) {
	wishData, err := json.Marshal(message.Data)
	if err != nil {
		log.Fatal("failed to marshal cars:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), wishData) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

func ValidationSearchResultHandler(message SearchResultMessage, topic int) {
	wishData, err := json.Marshal(message)
	if err != nil {
		log.Fatal("failed to marshal cars:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), wishData) // TODO Set key ?
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}
