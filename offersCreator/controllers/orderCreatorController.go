package controllers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	. "offersCreator/entities"
	"offersCreator/tools"
)

const STANDARD_CAR_PRICE = 100
var WISHES = make([]InitialWishRequest, 0)



func WishRequestedHandler(message InitialWishRequest) {
	// TODO save wish
	log.Println(message)
}

func RawWishHandler(rawWishesResult WishWithPossibilities, topic int) {
	enhanceOffer(&rawWishesResult.OfferPossibilities)

	//TODO do the scoring

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

func enhanceOffer(offersPossibilities *[]OfferPossibilities) {
	// Sort Offers by total amount
	//sort.Slice(offersPossibilities, func(i, j int) bool {
	//	return len((*offersPossibilities)[i].Offers) < len((*offersPossibilities)[j].Offers)
	//})

	//max := len((*offersPossibilities)[0].Offers)
	max := 200

	for _, offer := range *offersPossibilities {
		coefficient := float32(len(offer.Offers)) / float32(max)
		determinePrice(&offer, coefficient)
	}
}

func determinePrice(offerPossibilities *OfferPossibilities, coefficient float32) {
	var sum float32 = 0.0
	for _, offer := range offerPossibilities.Offers {
		offer.Price = (STANDARD_CAR_PRICE + float32(rand.Intn(10-(-10)) + (-10))) * coefficient*2
		sum += offer.Price
	}
	offerPossibilities.TotalPrice = sum
}
