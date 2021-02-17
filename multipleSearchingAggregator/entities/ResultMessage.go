package entities

type ResultMessage struct {
	OfferPossibilities	[]RawWishResult		`json:"offerPossibilities"`
	WishId				string				`json:"wishId"`
}
