package repository

import (
	"carBooking/entities"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var carType = entities.CarType{
	Name: "Liquid",
	Id:   1,
}
var carType2 = entities.CarType{
	Name: "Solid",
	Id:   2,
}

var node = entities.Node{
	Name:           "Nice",
	Id:             1,
	AvailableCarTypes: []entities.CarType{carType},
}
var node2 = entities.Node{
	Name:           "Marseilles",
	Id:             2,
	AvailableCarTypes: []entities.CarType{carType},
}
var nodes = []entities.Node{
	node, node2,
}

var car = entities.Car{
	Id: 1,
	CarType: carType,
}
var car2 = entities.Car{
	Id: 1,
	CarType: carType2,
}
var cars = []entities.Car{
	car, car2,
}

var bookings []entities.CarBooking

func InitMock(){
	date, _ := time.Parse(time.RFC3339, "2020-01-10T10:30:00+02:00")
	bookings = append(bookings, entities.CarBooking{
		Supplier: "Dracip",
		Date: date,
		Id: 1,
		Arrival: node2,
		Departure: node,
		Car: car,
	})

	date, _ = time.Parse(time.RFC3339, "2020-01-10T15:30:00+02:00")
	bookings = append(bookings, entities.CarBooking{
		Supplier: "Tahcapot",
		Date: date,
		Id: 2,
		Arrival: node2,
		Departure: node,
		Car: car2,
	})
}



func CreateBook(Date time.Time, Car entities.Car , Supplier string, NodeDeparture entities.Node, NodeArrival entities.Node){
	fmt.Println(Date)
	fmt.Println(Car)
	fmt.Println(Supplier)
	fmt.Println(NodeDeparture)
	fmt.Println(NodeArrival)

	bookings = append(bookings, entities.CarBooking{
		Date: Date,
		Supplier: Supplier,
		Departure: NodeDeparture,
		Arrival: NodeArrival,
		Car: Car,
	})
	//TODO use mongo
}

func FindAllBookings() []entities.CarBooking{
	return bookings
}

func GetNodeFromId(id int) (entities.Node, error){
	for _, node := range nodes{
		if node.Id == id{
			return node, nil
		}
	}
	return entities.Node{}, errors.New("Error 404: Node with id " + strconv.Itoa(id) + " not found")
}

func GetCarFromId(id int) (entities.Car, error){
	for _, car := range cars{
		if car.Id == id{
			return car, nil
		}
	}
	return entities.Car{}, errors.New("Error 404: Node with id " + strconv.Itoa(id) + " not found")
}
