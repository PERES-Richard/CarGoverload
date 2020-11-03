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
	assert.Equal(t, 1, res.Id)
	assert.Equal(t, 2, res.CarType.Id)
	assert.Equal(t, "Liquid", res.CarType.Name)

	res, err = repository.GetCarFromId(2)
	assert.Equal(t, nil, err)
	assert.IsType(t, entities.Car{}, res)
	assert.Equal(t, 2, res.Id)
	assert.Equal(t, 1, res.CarType.Id)
	assert.Equal(t, "Solid", res.CarType.Name)

	res, err = repository.GetCarFromId(-1)
	assert.NotEqual(t, nil, err)
	assert.IsType(t, entities.Car{}, res)
	assert.Equal(t, 0, res.Id)
}

func TestGetNodeFromId(t *testing.T) {
	var res, err = repository.GetNodeFromId(1)
	assert.Equal(t, nil, err)
	assert.IsType(t, entities.Node{}, res)
	assert.Equal(t, res.Id, 1)
	assert.Equal(t, res.Name, "Nice")

	res, err = repository.GetNodeFromId(2)
	assert.Equal(t, nil, err)
	assert.IsType(t, entities.Node{}, res)
	assert.Equal(t, 2, res.Id)
	assert.Equal(t, "Marseille", res.Name)

	res, err = repository.GetNodeFromId(-1)
	assert.NotEqual(t, nil, err)
	assert.IsType(t, entities.Node{}, res)
	assert.Equal(t, 0, res.Id)
}

func TestFindAll(t *testing.T){
	var res = repository.FindAllBookings(-1)
	assert.Equal(t, 1, res[0].Id)
	assert.Equal(t, 1, res[0].Car.Id)
	assert.Equal(t, 2, res[0].Car.CarType.Id)
	assert.Equal(t, "Liquid", res[0].Car.CarType.Name)
	assert.Equal(t, 1, res[0].Departure.Id)
	assert.Equal(t, "Nice", res[0].Departure.Name)
	assert.Equal(t, 2, res[0].Arrival.Id)
	assert.Equal(t, "Marseille", res[0].Arrival.Name)
	assert.Equal(t, "Picard", res[0].Supplier)

	assert.Equal(t, 2, res[1].Id)
	assert.Equal(t, "Amazoom", res[1].Supplier)
}

func TestFindAllTypeId(t *testing.T){
	var res = repository.FindAllBookings(1)
	assert.Equal(t, 2, res[0].Id)
	assert.Equal(t, "Amazoom", res[0].Supplier)
}

