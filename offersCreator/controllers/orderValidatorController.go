package controllers

import (
	"context"
	"encoding/json"
	"log"
	"offersCreator/tools"

	. "offersCreator/entities"
)

func WishRequestedHandler(message Wish) {
	// TODO save wish
}

func RawWishhandler(message Wish, topic int) {
	var isValid bool
	// TODO handler

	result, err := json.Marshal(isValid)
	if err != nil {
		log.Fatal("failed to marshal cars:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), result)
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}
