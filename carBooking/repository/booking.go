package repository

import "fmt"

func CreateBook(Date string, CarId int , Supplier string, NodeDepartureId int, NodeArrivalId int){
	fmt.Println(Date)
	fmt.Println(CarId)
	fmt.Println(Supplier)
	fmt.Println(NodeDepartureId)
	fmt.Println(NodeArrivalId)

	//TODO use mongo
}
