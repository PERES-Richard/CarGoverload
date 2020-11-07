package entities

import "encoding/json"

type Node struct {
	Name		string		`json:"name"`
	Id			json.Number `json:"id,number"`
	Latitude	float64		`json:"latitude"`
	Longitude	float64		`json:"longitude"`
}
