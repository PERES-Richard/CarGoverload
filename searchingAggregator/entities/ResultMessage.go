package entities

type ResultMessage struct {
	Offers		[]Offer		`json:"offers"`
	SearchId	string		`json:"searchId"`
}
