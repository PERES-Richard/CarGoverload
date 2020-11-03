package entities

type Car struct {
	Id 			int				`json:"id" pg:",pk"`
	CarTypeId	int				`json:"carTypeId"` //only used for pg orm
	CarType 	*CarType		`json:"carType" pg:"rel:has-one"`
}
