package entities

type SearchResultMessage struct {
	Offers 		[]Offer			`json:"offers"`
	SearchId 	string			`json:"searchId"`
}
