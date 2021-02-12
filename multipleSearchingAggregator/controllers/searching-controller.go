package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	. "multipleSearchingAggregator/entities"
	"multipleSearchingAggregator/tools"
)

const RAW_WISH_RESULT_TOPIC_WRITER_ID = 0

var searchArrayList []SearchData

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func SearchResultHandler(parsedMessage SearchResultMessage) {
	log.Println(parsedMessage)
	for i := range searchArrayList {
		for j := range searchArrayList[i].SearchIds {
			if searchArrayList[i].SearchIds[j] == parsedMessage.SearchId {
				//TODO remove doublons
				searchArrayList[i].SearchWithOffers[parsedMessage.SearchId] = parsedMessage.Offers
				searchArrayList[i].SearchesRemaining = searchArrayList[i].SearchesRemaining - 1
				if searchArrayList[i].SearchesRemaining == 0 {
					FinishAggregatingResults(searchArrayList[i])
				}
			}
		}
	}
}

func NewWishHandler(parsedMessage NewWishMessageResult) {
	searchArrayList = append(searchArrayList, SearchData{
		SearchIds: parsedMessage.SearchIds,
		WishId: parsedMessage.WishId,
		SearchesRemaining: len(parsedMessage.SearchIds),
		SearchWithOffers: make(map[string][]Offer),
	})
}

func FinishAggregatingResults(searchData SearchData) {
	log.Println(searchData)
	searchArrayList, _ = removeSearchData(searchData.WishId)

	offers := make([]Offer, 0)
	for _, value := range searchData.SearchWithOffers {
		for i := range value {
			offer := value[i]
			offers = append(offers, Offer{
				BookDate:  offer.BookDate,
				Arrival:   offer.Arrival,
				Departure: offer.Departure,
				Car:       offer.Car,
			})
		}
	}

	result := ResultMessage{
		Offers:   offers,
		WishId: searchData.WishId,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal("failed to marshal result:", err)
		return
	}

	log.Println("Results for wish ", searchData.WishId, " : ", offers)

	topic_id := RAW_WISH_RESULT_TOPIC_WRITER_ID
	kafkaErr := tools.KafkaPush(context.Background(), topic_id, []byte("value"), resultJSON)
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

func removeSearchData(wishId string) ([]SearchData, error) {
	err := errors.New("Remove error: car not found")
	var result []SearchData
	for _, s := range searchArrayList {
		if s.WishId != wishId {
			result = append(result, s)
		} else {
			err = errors.New("")
		}
	}
	return result, err
}
