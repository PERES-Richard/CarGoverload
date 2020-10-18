package entities

type Node struct {
	Name 						string		`json:"name"`
	Id 							int			`json:"id"`
	AvailableCarTypes []CarType				`json:"availableCarTypes"`
}
