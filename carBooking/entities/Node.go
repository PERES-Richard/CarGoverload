package entities

type Node struct {
	Name 						string		`json:"name"`
	Id 							int32			`json:"id"`
	AvailableCarTypes []CarType				`json:"availableCarTypes"`
}
