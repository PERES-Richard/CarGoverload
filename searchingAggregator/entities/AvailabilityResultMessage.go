package entities


type AvailabilityResultMessage struct {
	Cars 		[]int		`json:"cars"`
	SearchId 	int			`json:"searchId"`
}
