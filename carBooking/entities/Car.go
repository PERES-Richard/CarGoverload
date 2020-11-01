package entities

type Car struct {
	Id 			int64			`json:"id"`
	CarTypeId	int64			`json:"carTypeId"` //only used for pg orm
	CarType 	*CarType		`json:"carType" pg:"rel:has-one"`
}
