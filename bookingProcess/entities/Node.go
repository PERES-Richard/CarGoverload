package entities

type Node struct {
	Name 						string		`json:"name"`
	Id 							int			`json:"id"`
	Latitude	float64		`json:"latitude"`
	Longitude	float64		`json:"longitude"`
}
