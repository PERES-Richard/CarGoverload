package model

type Node struct {
	Name 						string		`json:"name"`
	Id 							int			`json:"id"`
	AvailableTypes []CarType				`json:"availableTypes"`
}
