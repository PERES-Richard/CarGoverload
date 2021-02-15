package controllers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"offersCreator/tools"
	"sort"

	. "offersCreator/entities"
)

const STANDARD_CAR_PRICE = 100
var WISHES = make([]Wish, 0)



func WishRequestedHandler(message Wish) {
	// TODO save wish
}

func RawWishHandler(rawWishesResult []Book, topic int) {
	enhanceBook(&rawWishesResult)

	offers, err := json.Marshal(rawWishesResult)
	if err != nil {
		log.Fatal("failed to marshal offers:", err)
		return
	}

	kafkaErr := tools.KafkaPush(context.Background(), topic, []byte("value"), offers)
	if kafkaErr != nil {
		log.Panic("failed to write message:", kafkaErr)
	}
}

func enhanceBook(rawWishesResult *[]Book) {
	// Sort Offers by total amount
	sort.Slice(rawWishesResult, func(i, j int) bool {
		return len((*rawWishesResult)[i].Offers) < len((*rawWishesResult)[j].Offers)
	})

	max := len((*rawWishesResult)[0].Offers)

	for _, book := range *rawWishesResult {
		coeff := float32(len(book.Offers)) / float32(max)
		determinePrice(&book, coeff)
	}
}

func determinePrice(book *Book, coeff float32) {
	var sum float32 = 0.0
	for _, offer := range book.Offers {
		offer.Price = (STANDARD_CAR_PRICE + float32(rand.Intn(10-(-10)) + (-10))) * coeff*2
		sum += offer.Price
	}
	book.TotalPrice = sum
}
