package entities


type AvailabilityResultMessage struct {
	Cars 		[]int		`json:"cars"`
	SearchId 	string			`json:"searchId"`
}
