package repository

import (
	"carBooking/entities"
	"database/sql"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

var dbHost string
var dbPort int
var dbPassword string
var dbName string
var dbUser string

var InTest = false

var db *pg.DB

func initDatabaseClient(){
	db = pg.Connect(&pg.Options{
		User: dbUser,
		Password: dbPassword,
		Database: dbName,
	})
}

func initTestDatabase(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		panic(err)
	}
}

func clearDatabase(){
	//var CollectionNode = "node"
	//var CollectionBooking = "booking"
	//var CollectionCarType = "car_type"
	//var CollectionNodeCarType = "node_car_type"
	//var CollectionCar = "car"
	//
	//_, err := db.Exec("DROP TABLE IF EXISTS " + CollectionBooking)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionNodeCarType)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionNode)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionCar)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionCarType)
	//if err != nil {
	//	panic(err)
	//}
}

func createTables(){
	models := []interface{}{
		(*entities.CarType)(nil),
		(*entities.Car)(nil),
		(*entities.Node)(nil),
		(*entities.CarBooking)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			panic(err)
		}
	}
}

func populateTables(){
	// MOCKING CarType
	carTypeSolid := &entities.CarType{
		Name:   "Solid",
	}
	_, err := db.Model(carTypeSolid).Insert()
	if err != nil {
		panic(err)
	}
	carTypeLiquid := &entities.CarType{
		Name:   "Liquid",
	}
	_, err = db.Model(carTypeLiquid).Insert()
	if err != nil {
		panic(err)
	}
	// MOCKING Car
	car1 := &entities.Car{
		CarTypeId: carTypeLiquid.Id,
	}
	_, err = db.Model(car1).Insert()
	if err != nil {
		panic(err)
	}
	car2 := &entities.Car{
		CarTypeId: carTypeSolid.Id,
	}
	_, err = db.Model(car2).Insert()
	if err != nil {
		panic(err)
	}
	car3 := &entities.Car{
		CarTypeId: carTypeSolid.Id,
	}
	_, err = db.Model(car3).Insert()
	if err != nil {
		panic(err)
	}
	car4 := &entities.Car{
		CarTypeId: carTypeSolid.Id,
	}
	_, err = db.Model(car4).Insert()
	if err != nil {
		panic(err)
	}
	car5 := &entities.Car{
		CarTypeId: carTypeLiquid.Id,
	}
	_, err = db.Model(car5).Insert()
	if err != nil {
		panic(err)
	}
	// MOCKING node
	nodeNice := &entities.Node{
		Name: "Nice",
		AvailableCarTypes: []entities.CarType{
			*carTypeLiquid,
		},
	}
	_, err = db.Model(nodeNice).Insert()
	if err != nil {
		panic(err)
	}
	nodeMarseille := &entities.Node{
		Name: "Marseille",
		AvailableCarTypes: []entities.CarType{
			*carTypeLiquid,
			*carTypeSolid,
		},
	}
	_, err = db.Model(nodeMarseille).Insert()
	if err != nil {
		panic(err)
	}
	nodeDraguignan := &entities.Node{
		Name: "Draguignan",
		AvailableCarTypes: []entities.CarType{
			*carTypeLiquid,
			*carTypeSolid,
		},
	}
	_, err = db.Model(nodeDraguignan).Insert()
	if err != nil {
		panic(err)
	}
	nodeToulon := &entities.Node{
		Name: "Toulon",
		AvailableCarTypes: []entities.CarType{
			*carTypeSolid,
		},
	}
	_, err = db.Model(nodeToulon).Insert()
	if err != nil {
		panic(err)
	}
	nodeLyon := &entities.Node{
		Name: "Lyon",
		AvailableCarTypes: []entities.CarType{
			*carTypeSolid,
		},
	}
	_, err = db.Model(nodeLyon).Insert()
	if err != nil {
		panic(err)
	}
	nodeAvignon := &entities.Node{
		Name: "Avignon",
		AvailableCarTypes: []entities.CarType{
			*carTypeLiquid,
		},
	}
	_, err = db.Model(nodeAvignon).Insert()
	if err != nil {
		panic(err)
	}
	nodeParis := &entities.Node{
		Name: "Paris",
		AvailableCarTypes: []entities.CarType{
			*carTypeLiquid,
			*carTypeSolid,
		},
	}
	_, err = db.Model(nodeParis).Insert()
	if err != nil {
		panic(err)
	}

	car := new(entities.Car)
	err = db.Model(car).
		Relation("CarType").
		Where("car.id = ?", car1.Id).
		Select()
	if err != nil {
		panic(err)
	}

	carX1, _ := GetCarFromId(1)
	carX2, _ := GetCarFromId(2)
	carX3, _ := GetCarFromId(3)
	nodeNiceX, _ := GetNodeFromId(1)
	nodeMarseilleX, _ := GetNodeFromId(2)
	nodeDraguignanX, _ := GetNodeFromId(3)
	nodeToulonX, _ := GetNodeFromId(4)
	nodeLyonX, _ := GetNodeFromId(5)
	CreateBook(time.Now(), &carX1, "Picard", &nodeNiceX, &nodeMarseilleX)
	CreateBook(time.Now().Add(5000), &carX2, "Amazoom", &nodeDraguignanX, &nodeMarseilleX)
	CreateBook(time.Now().Add(50000), &carX2, "Microsoft", &nodeToulonX, &nodeNiceX)
	CreateBook(time.Now().Add(50000), &carX3, "Fnac", &nodeMarseilleX, &nodeNiceX)
	CreateBook(time.Now().Add(100000), &carX3, "Darty", &nodeLyonX, &nodeToulonX)
}

func InitDatabase(){
	dbHost = os.Getenv("DB_HOST")
	dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")

	if InTest{
		if dbHost == ""{ //not launched with docker
			dbHost = "localhost"
			dbPort = 5432
			dbPassword = "superpassword"
			dbUser = "cargoverload"
		}
		dbName = "cargoverload_test"
		initTestDatabase()
	}

	initDatabaseClient()
	clearDatabase()
	createTables()
	populateTables()
}

func CreateBook(date time.Time, car *entities.Car , supplier string, nodeDeparture *entities.Node, nodeArrival *entities.Node) entities.CarBooking{
	var booking = &entities.CarBooking{
		Date: date,
		Supplier: supplier,
		DepartureId: nodeDeparture.Id,
		Departure: nodeDeparture,
		ArrivalId: nodeDeparture.Id,
		Arrival: nodeArrival,
		CarId: car.Id,
		Car: car,
	}

	_, err := db.Model(booking).Insert()
	if err != nil {
		panic(err)
	}
	return *booking
}

func GetAllNodes() []entities.Node{
	var nodes =  []entities.Node{}
	err := db.Model(&nodes).Relation("AvailableCarTypes").Select()
	if err != nil {
		panic(err)
	}
	return nodes
}

func GetAllTypes() []entities.CarType{
	var carTypes = []entities.CarType{}
	err := db.Model(&carTypes).Select()
	if err != nil {
		panic(err)
	}
	return carTypes
}

func FindAllBookings(typeId int64) []entities.CarBooking{
	var bookings []entities.CarBooking
	err := db.Model(&bookings).
		Relation("Arrival").
		Relation("Departure").
		Relation("Car").
		Relation("Car.CarType").
		Relation("Departure.AvailableCarTypes").
		Relation("Arrival.AvailableCarTypes").
		Select()
	if err != nil {
		panic(err)
	}

	var typedBooking = []entities.CarBooking{} //uggly way because orm is bad with WHERE and not working
	if typeId != -1{
		for _, book := range bookings {
			if book.Car.CarType.Id == typeId{
				typedBooking = append(typedBooking, book)
			}
		}
	}else{
		typedBooking = bookings
	}

	//TODO add available car types
	//for _, book := range bookings {
	//	book.Departure.AvailableCarTypes = GetCarTypesFromId(book.DepartureId)
	//	book.Arrival.AvailableCarTypes = GetCarTypesFromId(book.ArrivalId)
	//}

	return typedBooking
}

func GetNodeFromId(id int) (entities.Node, error){
	node := new(entities.Node)
	err := db.Model(node).
		Relation("AvailableCarTypes").
		Where("node.id = ?", id).
		Select()
	if err != nil {
		panic(err)
	}
	//err = db.Close()
	//if err != nil {
	//	panic(err)
	//}
	return *node, nil
}

func GetCarFromId(id int) (entities.Car, error) {
	car := new(entities.Car)
	err := db.Model(car).
		Relation("CarType").
		Where("car.id = ?", id).
		Select()
	if err != nil {
		panic(err)
	}
	//err = db.Close()
	//if err != nil {
	//	panic(err)
	//}
	return *car, nil
}
