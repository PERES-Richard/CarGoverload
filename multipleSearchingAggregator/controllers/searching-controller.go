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
	log.Println("Search received for seardhId : ", parsedMessage.SearchId)
	for i := range searchArrayList {
		for j := range searchArrayList[i].SearchIds {
			if searchArrayList[i].SearchIds[j] == parsedMessage.SearchId {
				searchArrayList[i].SearchWithOffers[parsedMessage.SearchId] = parsedMessage.Offers
				searchArrayList[i].SearchesRemaining = searchArrayList[i].SearchesRemaining - 1
				log.Println("Test nombre search restantes : ", searchArrayList[i].SearchesRemaining)
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
	log.Println("Entering finish aggregate")
	removeDuplicates(searchData)

	searchArrayList, _ = removeSearchData(searchData.WishId)

	rawWishResults := make([]RawWishResult, 0)
	for key, value := range searchData.SearchWithOffers {
		offers := make([]Offer, 0)
		for i := range value {
			offer := value[i]
			offers = append(offers, Offer{
				BookDate:  offer.BookDate,
				Arrival:   offer.Arrival,
				Departure: offer.Departure,
				Car:       offer.Car,
				Distance:  offer.Distance,
			})
		}
		rawWishResults = append(rawWishResults, RawWishResult{
			SearchId: key,
			Offers: offers,
		})
	}

	result := ResultMessage{
		OfferPossibilities:  rawWishResults,
		WishId: searchData.WishId,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal("failed to marshal result:", err)
		return
	}

	log.Println("Results for wish ", searchData.WishId, " : ", rawWishResults)

	kafkaErr := tools.KafkaPush(context.Background(), RAW_WISH_RESULT_TOPIC_WRITER_ID, []byte("value"), resultJSON)
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

func removeDuplicates(searchData SearchData) {
	offersInWish :=  make(map[string][]Offer)
	carsAlreadyUsed := make([]Car, 0)
	for key, value := range searchData.SearchWithOffers {
		carsOfOffer := make([]Offer, 0)
		for i := range value {
			carsOfOffer = append(carsOfOffer, value[i])
		}
		offersInWish[key] = carsOfOffer
	}
	result :=  make(map[string][]Offer)
	for isARemainingOffer(offersInWish) {
		for key, value := range offersInWish {
			offersInWish[key] = removeOffersWithCarsAlreadyUsed(value, carsAlreadyUsed)
			if len(offersInWish[key]) > 0 { // on vérifie si la search a toujours des offers disponibles
				offerToKeep := offersInWish[key][0]
				result[key] = append(result[key], offerToKeep)                              // on ajoute au résultat la première offer pour cette recherche puis on passe à la suivante etc
				carsAlreadyUsed = append(carsAlreadyUsed, offerToKeep.Car)                  // la car est maintenant use dans une offer
				offersInWish[key] = removeFromArrayAtIndex(0, offersInWish[key]) // remove la premiere offer car maintenant on l'a use
			}
		}
	}
	for key, value := range result {
		searchData.SearchWithOffers[key] = value
	}
	log.Println("Final result : ", result)
}

func removeFromArrayAtIndex(index int, offers []Offer) []Offer {
	log.Println("Before removing : ", offers)
	offers = append(offers[:index], offers[index + 1:]...)
	log.Println("After removing : ", offers)
	return offers
}

func isARemainingOffer(offersInWish map[string][]Offer) bool {
	for key := range offersInWish {
		if len(offersInWish[key]) > 0 {
			return true
		}
	}
	return false
}

func removeOffersWithCarsAlreadyUsed(offers []Offer, cars []Car) []Offer {
	log.Println("Entering first remove : ", offers)
	result := make([]Offer, 0)
	for i := range offers {
		log.Println("Remove offers : ", offers[i])
		offer := offers[i]
		if !isCarAlreadyUsed(offer.Car, cars) {
			result = append(result, offer)
		}
	}
	return result
}

func isCarAlreadyUsed(car Car, cars []Car) bool {
	for i := range cars {
		if cars[i].Id == car.Id {
			return true
		}
	}
	return false
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
