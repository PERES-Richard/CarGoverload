package entities

type LocationResultMessage struct {
	SearchId string       `json:"searchId"`
	Cars     []TrackedCar `json:"results"`
}
