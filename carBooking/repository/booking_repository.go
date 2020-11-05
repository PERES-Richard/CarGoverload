package repository

import (
	"carBooking/entities"
	"database/sql"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	_ "github.com/lib/pq"
	"log"
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
		log.Fatal(err)
	}
	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		log.Fatal(err)
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
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionNodeCarType)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionNode)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionCar)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionCarType)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func createTables(){
	models := []interface{}{
		(*entities.Car)(nil),
		(*entities.Node)(nil),
		(*entities.CarBooking)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func populateTables(){
	// MOCKING Car
	car0 := &entities.Car{
		CarTypeId: 1,
		Id: 1,
	}
	_, err := db.Model(car0).Insert()
	if err != nil {
		log.Fatal(err)
	}
	car1 := &entities.Car{
		CarTypeId: 2,
		Id: 2,
	}
	_, err = db.Model(car1).Insert()
	if err != nil {
		log.Fatal(err)
	}
	car2 := &entities.Car{
		CarTypeId: 2,
		Id: 3,
	}
	_, err = db.Model(car2).Insert()
	if err != nil {
		log.Fatal(err)
	}
	car3 := &entities.Car{
		CarTypeId: 1,
		Id: 4,
	}
	_, err = db.Model(car3).Insert()
	if err != nil {
		log.Fatal(err)
	}
	car4 := &entities.Car{
		CarTypeId: 2,
		Id: 5,
	}
	_, err = db.Model(car4).Insert()
	if err != nil {
		log.Fatal(err)
	}
	// MOCKING node
	node0 := &entities.Node{
		Name: "Marseille",
		Id: 1,
	}
	_, err = db.Model(node0).Insert()
	if err != nil {
		log.Fatal(err)
	}
	node1 := &entities.Node{
		Name: "Avignon-liquid",
		Id: 2,
	}
	_, err = db.Model(node1).Insert()
	if err != nil {
		log.Fatal(err)
	}
	node2 := &entities.Node{
		Name: "Nice",
		Id: 3,
	}
	_, err = db.Model(node2).Insert()
	if err != nil {
		log.Fatal(err)
	}
	node3 := &entities.Node{
		Name: "Paris",
		Id: 4,
	}
	_, err = db.Model(node3).Insert()
	if err != nil {
		log.Fatal(err)
	}
	node4 := &entities.Node{
		Name: "Avignon-solid",
		Id: 5,
	}
	_, err = db.Model(node4).Insert()
	if err != nil {
		log.Fatal(err)
	}

	carX1, _ := GetCarFromId(1)
	carX2, _ := GetCarFromId(2)
	carX3, _ := GetCarFromId(3)
	carX4, _ := GetCarFromId(4)
	carX5, _ := GetCarFromId(5)
	node0X, _ := GetNodeFromId(1)
	node1X, _ := GetNodeFromId(2)
	node2X, _ := GetNodeFromId(3)
	node3X, _ := GetNodeFromId(4)
	node4X, _ := GetNodeFromId(5)
	CreateBook(time.Now().Add(time.Hour * time.Duration(0)), time.Now().Add(time.Hour * time.Duration(1)), &carX1, "Picard", &node0X, &node1X)
	CreateBook(time.Now().Add(time.Hour * time.Duration(0)), time.Now().Add(time.Hour * time.Duration(1)), &carX1, "Amazoom", &node2X, &node3X)
	CreateBook(time.Now().Add(time.Hour * time.Duration(0)), time.Now().Add(time.Hour * time.Duration(1)), &carX2, "Microsoft", &node0X, &node4X)
	CreateBook(time.Now().Add(time.Hour * time.Duration(0)), time.Now().Add(time.Hour * time.Duration(1)), &carX1, "Fnac", &node3X, &node2X)
	CreateBook(time.Now().Add(time.Hour * time.Duration(1)), time.Now().Add(time.Hour * time.Duration(2)), &carX3, "Darty", &node4X, &node0X)
	CreateBook(time.Now().Add(time.Hour * time.Duration(2)), time.Now().Add(time.Hour * time.Duration(3)), &carX4, "TopAchat", &node1X, &node2X)
	CreateBook(time.Now().Add(time.Hour * time.Duration(3)), time.Now().Add(time.Hour * time.Duration(4)), &carX5, "LDLC", &node3X, &node1X)
	CreateBook(time.Now().Add(time.Hour * time.Duration(4)), time.Now().Add(time.Hour * time.Duration(5)), &carX3, "MiamMiam", &node0X, &node2X)
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

func CreateBook(date time.Time, dateArrival time.Time, car *entities.Car , supplier string, nodeDeparture *entities.Node, nodeArrival *entities.Node) entities.CarBooking{
	var booking = &entities.CarBooking{
		Date: date,
		Supplier: supplier,
		DepartureId: nodeDeparture.Id,
		Departure: nodeDeparture,
		ArrivalId: nodeArrival.Id,
		Arrival: nodeArrival,
		CarId: car.Id,
		Car: car,
		DateArrival: dateArrival,
	}

	_, err := db.Model(booking).Insert()
	if err != nil {
		log.Fatal(err)
	}
	return *booking
}

func FindAllBookings(typeId int) []entities.CarBooking{
	var bookings []entities.CarBooking
	err := db.Model(&bookings).
		Relation("Arrival").
		Relation("Departure").
		Relation("Car").
		Select()
	if err != nil {
		log.Fatal(err)
	}

	var typedBooking = []entities.CarBooking{} //uggly way because orm is bad with WHERE and not working
	if typeId != -1{
		for _, book := range bookings {
			if book.Car.CarTypeId == typeId{
				typedBooking = append(typedBooking, book)
			}
		}
	}else{
		typedBooking = bookings
	}

	return typedBooking
}

func GetNodeFromId(id int) (entities.Node, error){
	node := new(entities.Node)
	err := db.Model(node).
		Where("node.id = ?", id).
		Select()
	if err != nil {
		return entities.Node{}, err
	}
	//err = db.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
	return *node, nil
}

func GetCarFromId(id int) (entities.Car, error) {
	car := new(entities.Car)
	err := db.Model(car).
		Where("car.id = ?", id).
		Select()
	if err != nil {
		return entities.Car{}, err
	}
	//err = db.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
	return *car, nil
}
