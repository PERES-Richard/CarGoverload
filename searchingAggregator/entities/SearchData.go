package entities

import "time"

type SearchData struct {
	LocationResult			*[]TrackedCar
	AvailabilityResult 		*[]int
	SearchTime				*time.Time
	Validation				bool
}
