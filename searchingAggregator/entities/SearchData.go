package entities

import "time"

type SearchData struct {
	SearchId				string
	LocationResult			[]TrackedCar
	AvailabilityResult 		[]int
	SearchTime				time.Time
	Validation				bool
	ReceivedLocation		bool
	ReceivedAvailability	bool
	ReceivedTime			bool
}
