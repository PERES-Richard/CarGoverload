package entities

type LocationResultMessage struct {
	SearchId	int				`json:"searchId"`
	Cars 		[]TrackedCar	`json:"results"`
}