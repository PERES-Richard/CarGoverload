package entities

type RawWishResult struct {
	SearchId 		string 		`json:"searchId"`
	Offers	 		[]Offer		`json:"offers"`
}
