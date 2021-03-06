package entities


type AvailabilityResultMessage struct {
	Cars 		[]int			`json:"carIdsBooked"`
	SearchId 	string			`json:"searchId"`
}
