package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCarAvailability(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CarAvailability Suite")
}
