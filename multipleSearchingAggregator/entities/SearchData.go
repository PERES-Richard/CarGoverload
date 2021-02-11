package entities

type SearchData struct {
	SearchIds				[]string
	SearchWithOffers		map[string][]Offer
	SearchesRemaining		int
	WishId					string
}
