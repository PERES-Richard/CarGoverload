package entities

type Supplier struct {
	ID        int      `json:"id"`
	Name     string    `json:"name"`
	Offers	 []Offer `json:"offers"`
}
