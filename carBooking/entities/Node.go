package entities

type Node struct {
	Name 						string		`json:"name"`
	Id 							int64		`json:"id"`
	AvailableCarTypes 			[]CarType	`json:"availableCarTypes" pg:"rel:has-many"`
}
