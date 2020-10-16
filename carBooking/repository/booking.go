package repository

import (
	"carBooking/model"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var carType = model.CarType{
	Name: "Liquid",
	Id:   1,
}
var carType2 = model.CarType{
	Name: "Solid",
	Id:   2,
}

var node = model.Node{
	Name:           "Nice",
	Id:             1,
	AvailableCarTypes: []model.CarType{carType},
}
var node2 = model.Node{
	Name:           "Marseilles",
	Id:             2,
	AvailableCarTypes: []model.CarType{carType},
}
var nodes = []model.Node{
	node, node2,
}

var car = model.Car{
	Id: 1,
	CarType: carType,
}
var car2 = model.Car{
	Id: 1,
	CarType: carType2,
}
var cars = []model.Car{
	car, car2,
}

var bookings []model.CarBooking

func InitMock(){
	date, _ := time.Parse(time.RFC3339, "2020-01-10T10:30:00+02:00")
	bookings = append(bookings, model.CarBooking{
		Supplier: "Dracip",
		Date: date,
		Id: 1,
		Arrival: node2,
		Departure: node,
		Car: car,
	})

	date, _ = time.Parse(time.RFC3339, "2020-01-10T15:30:00+02:00")
	bookings = append(bookings, model.CarBooking{
		Supplier: "Tahcapot",
		Date: date,
		Id: 2,
		Arrival: node2,
		Departure: node,
		Car: car2,
	})
}



func CreateBook(Date time.Time, Car model.Car , Supplier string, NodeDeparture model.Node, NodeArrival model.Node){
	fmt.Println(Date)
	fmt.Println(Car)
	fmt.Println(Supplier)
	fmt.Println(NodeDeparture)
	fmt.Println(NodeArrival)

	bookings = append(bookings, model.CarBooking{
		Date: Date,
		Supplier: Supplier,
		Departure: NodeDeparture,
		Arrival: NodeArrival,
		Car: Car,
	})
	//TODO use mongo
}

func FindAllBookings() []model.CarBooking{
	return bookings
}

func GetNodeFromId(id int) (model.Node, error){
	for _, node := range nodes{
		if node.Id == id{
			return node, nil
		}
	}
	return model.Node{}, errors.New("Error 404: Node with id " + strconv.Itoa(id) + " not found")
}

func GetCarFromId(id int) (model.Car, error){
	for _, car := range cars{
		if car.Id == id{
			return car, nil
		}
	}
	return model.Car{}, errors.New("Error 404: Node with id " + strconv.Itoa(id) + " not found")
}
