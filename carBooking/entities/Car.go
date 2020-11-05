package entities

type Car struct {
	Id 			int				`json:"id" pg:",pk"`
	CarTypeId	int				`json:"carTypeId"`
}
