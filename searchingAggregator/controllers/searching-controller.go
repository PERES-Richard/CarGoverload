package controllers

import (
	. "searchingAggregator/entities"
)

// Custom error to return in case of a JSON parsing error
type JSONError struct {
	Message string `json:"Message"`
}

func AvailabilityResultHandler(parsedMessage AvailabilityResultMessage) {

}

func LocationResultHandler(parsedMessage LocationResultMessage) {

}

func NewSearchHandler(parsedMessage SearchMessage) {

}

func NewValidationSearchHandler(parsedMessage SearchMessage) {
	NewSearchHandler(parsedMessage)
}