package entities

type ResultMessage struct {
	Offers			[]Offer			`json:"offers"`
	WishId			string			`json:"wishId"`
}
