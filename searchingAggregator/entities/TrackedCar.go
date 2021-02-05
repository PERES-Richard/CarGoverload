package entities

type TrackedCar struct {
	Car			Car			`json:"car"`
	Node 		Node		`json:"node"`
	DestNode	Node		`json:"destinationNode"`
}
