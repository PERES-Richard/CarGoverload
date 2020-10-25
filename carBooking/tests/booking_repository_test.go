package tests

import (
	"carBooking/entities"
	"carBooking/repository"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setUp(){
	repository.InTest = true
	repository.InitDatabase()
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestGetCarFromId(t *testing.T) {
	var res, err = repository.GetCarFromId(1)
	assert.Equal(t, nil, err)
	assert.IsType(t, entities.Car{}, res)
	assert.Equal(t, res.Id, 1)
	assert.Equal(t, res.CarType.Id, 2)
	assert.Equal(t, res.CarType.Name, "Liquid")

	res, err = repository.GetCarFromId(2)
	assert.Equal(t, nil, err)
	assert.IsType(t, entities.Car{}, res)
	assert.Equal(t, res.Id, 2)
	assert.Equal(t, res.CarType.Id, 1)
	assert.Equal(t, res.CarType.Name, "Solid")

	res, err = repository.GetCarFromId(-1)
	assert.NotEqual(t, nil, err)
	assert.IsType(t, entities.Car{}, res)
	assert.Equal(t, res.Id, 0)
	assert.Equal(t, res.CarType.Id, 0)
	assert.Equal(t, res.CarType.Name, "")
}

func TestGetCarTypesForNode(t *testing.T){
	var res = repository.GetCarTypesForNode(1)
	assert.Equal(t, res[0].Id, 2)
	assert.Equal(t, res[0].Name, "Liquid")

	res = repository.GetCarTypesForNode(2)
	assert.Equal(t, res[0].Id, 1)
	assert.Equal(t, res[0].Name, "Solid")
	assert.Equal(t, res[1].Id, 2)
	assert.Equal(t, res[1].Name, "Liquid")
}

func TestGetNodeFromId(t *testing.T) {
	var res, err = repository.GetNodeFromId(1)
	assert.Equal(t, nil, err)
	assert.IsType(t, entities.Node{}, res)
	assert.Equal(t, res.Id, 1)
	assert.Equal(t, res.Name, "Nice")
	assert.Equal(t, res.AvailableCarTypes[0].Id, 2)
	assert.Equal(t, res.AvailableCarTypes[0].Name, "Liquid")

	res, err = repository.GetNodeFromId(2)
	assert.Equal(t, nil, err)
	assert.IsType(t, entities.Node{}, res)
	assert.Equal(t, res.Id, 2)
	assert.Equal(t, res.Name, "Marseille")
	assert.Equal(t, res.AvailableCarTypes[0].Id, 1)
	assert.Equal(t, res.AvailableCarTypes[0].Name, "Solid")
	assert.Equal(t, res.AvailableCarTypes[1].Id, 2)
	assert.Equal(t, res.AvailableCarTypes[1].Name, "Liquid")

	res, err = repository.GetNodeFromId(-1)
	assert.NotEqual(t, nil, err)
	assert.IsType(t, entities.Node{}, res)
	assert.Equal(t, res.Id, 0)
	assert.Equal(t, res.Name, "")
}

func TestFindAll(t *testing.T){
	var res = repository.FindAllBookings(-1)
	assert.Equal(t, res[0].Id, 1)
	assert.Equal(t, res[0].Car.Id, 1)
	assert.Equal(t, res[0].Car.CarType.Id, 2)
	assert.Equal(t, res[0].Car.CarType.Name, "Liquid")
	assert.Equal(t, res[0].Departure.Id, 1)
	assert.Equal(t, res[0].Departure.Name, "Nice")
	assert.Equal(t, res[0].Departure.AvailableCarTypes[0].Id, 2)
	assert.Equal(t, res[0].Departure.AvailableCarTypes[0].Name, "Liquid")
	assert.Equal(t, res[0].Arrival.Id, 2)
	assert.Equal(t, res[0].Arrival.Name, "Marseille")
	assert.Equal(t, res[0].Arrival.AvailableCarTypes[0].Id, 1)
	assert.Equal(t, res[0].Arrival.AvailableCarTypes[1].Id, 2)
	assert.Equal(t, res[0].Arrival.AvailableCarTypes[0].Name, "Solid")
	assert.Equal(t, res[0].Arrival.AvailableCarTypes[1].Name, "Liquid")
	assert.Equal(t, res[0].Supplier, "Picard")

	assert.Equal(t, res[1].Id, 2)
	assert.Equal(t, res[1].Supplier, "Amazoom")
}

func TestFindAllTypeId(t *testing.T){
	var res = repository.FindAllBookings(1)
	assert.Equal(t, res[0].Id, 2)
	assert.Equal(t, res[0].Supplier, "Amazoom")
}

