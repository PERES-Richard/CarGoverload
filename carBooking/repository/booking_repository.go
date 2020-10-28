package repository

import (
	"carBooking/entities"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

var CollectionNode = "node"
var CollectionBooking = "booking"
var CollectionCarType = "car_type"
var CollectionNodeCarType = "node_car_type"
var CollectionCar = "car"

var dbHost string
var dbPort int
var dbPassword string
var dbName string
var dbUser string

var InTest = false

func getDatabaseClient() *sql.DB{
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//DO not forget defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
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

func clearDatabase(db *sql.DB){
	_, err := db.Exec("DROP TABLE IF EXISTS " + CollectionBooking)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionNodeCarType)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionNode)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionCar)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS " + CollectionCarType)
	if err != nil {
		panic(err)
	}
}

func createTables(db *sql.DB){
	_, err := db.Exec(`CREATE TABLE ` + CollectionCarType + ` (
    							id SERIAL PRIMARY KEY UNIQUE, 
    							name VARCHAR(100)
						   );`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE ` + CollectionCar + ` (
    							id SERIAL PRIMARY KEY UNIQUE, 
    							car_type integer,
    							FOREIGN KEY(car_type) REFERENCES car_type(id)
						   );`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE ` + CollectionNode + ` (
    							id SERIAL PRIMARY KEY UNIQUE, 
    							name VARCHAR(100)
						   );`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE ` + CollectionNodeCarType +` (
    							node_id integer, 
    							car_type_id integer,
    							PRIMARY KEY (node_id, car_type_id),
    							FOREIGN KEY (node_id) REFERENCES node(id),
    							FOREIGN KEY (car_type_id) REFERENCES car_type(id)
						   );`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE ` + CollectionBooking + ` (
    							id SERIAL PRIMARY KEY UNIQUE, 
    							car_id integer,
    							supplier VARCHAR(100),
    							departure_id integer,
    							arrival_id integer,
								time varchar(100),
    							FOREIGN KEY (departure_id) REFERENCES node(id),
    							FOREIGN KEY (arrival_id) REFERENCES node(id)
						   );`)
	if err != nil {
		panic(err)
	}
}

func populateTables(db *sql.DB){
	_, err := db.Exec(`INSERT INTO ` + CollectionCarType + ` (name) VALUES($1)`, "Solid")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionCarType + ` (name) VALUES($1)`, "Liquid")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO ` + CollectionCar + ` (car_type) VALUES($1)`, 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionCar + ` (car_type) VALUES($1)`, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionCar + ` (car_type) VALUES($1)`, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionCar + ` (car_type) VALUES($1)`, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionCar + ` (car_type) VALUES($1)`, 2)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO ` + CollectionNode + ` (name) VALUES($1)`, "Nice")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNode + ` (name) VALUES($1)`, "Marseille")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNode + ` (name) VALUES($1)`, "Draguignan")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNode + ` (name) VALUES($1)`, "Toulon")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNode + ` (name) VALUES($1)`, "Lyon")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNode + ` (name) VALUES($1)`, "Avignon")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNode + ` (name) VALUES($1)`, "Paris")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 1, 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 2, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 2, 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 3, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 3, 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 4, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 5, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 6, 2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 7, 1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO ` + CollectionNodeCarType + ` (node_id, car_type_id) VALUES($1, $2)`, 7, 2)
	if err != nil {
		panic(err)
	}

	car1, _ := GetCarFromId(1)
	car2, _ := GetCarFromId(2)
	nodeNice, _ := GetNodeFromId(1)
	nodeMarseille, _ := GetNodeFromId(2)
	nodeDraguignan, _ := GetNodeFromId(3)
	CreateBook(time.Now(), car1, "Picard", nodeNice, nodeMarseille)
	CreateBook(time.Now(), car2, "Amazoom", nodeDraguignan, nodeMarseille)
	CreateBook(time.Now(), car2, "Test", nodeDraguignan, nodeMarseille)
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

	db := getDatabaseClient()
	clearDatabase(db)
	createTables(db)
	populateTables(db)
}

func CreateBook(date time.Time, car entities.Car , supplier string, nodeDeparture entities.Node, nodeArrival entities.Node) entities.CarBooking{
	db := getDatabaseClient()

	var booking = entities.CarBooking{
		Date: date,
		Supplier: supplier,
		Departure: nodeDeparture,
		Arrival: nodeArrival,
		Car: car,
	}

	_, err := db.Exec(`INSERT INTO ` + CollectionBooking + ` (car_id, supplier, departure_id, arrival_id, time) VALUES($1, $2, $3, $4, $5)`, car.Id, supplier, nodeDeparture.Id, nodeArrival.Id, date.Unix())
	if err != nil {
		panic(err)
	}

	return booking
}

func GetAllNodes() []entities.Node{
	var nodes = []entities.Node{}
	db := getDatabaseClient()

	var rows *sql.Rows
	var err error
	rows, err = db.Query(`SELECT id, name FROM ` + CollectionNode + ` ORDER BY name`)
	if err != nil {
		_ = db.Close()
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var nodeId int
		var nodeName string
		err = rows.Scan(&nodeId, &nodeName)
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, entities.Node{
			Name:              nodeName,
			Id:                nodeId,
			AvailableCarTypes: GetCarTypesForNode(nodeId),
		})
	}
	err = rows.Err()
	if err != nil {
		_ = db.Close()
		panic(err)
	}
	_ = db.Close()
	return nodes
}

func GetAllTypes() []entities.CarType{
	var nodes = []entities.CarType{}
	db := getDatabaseClient()

	var rows *sql.Rows
	var err error
	rows, err = db.Query(`SELECT id, name FROM ` + CollectionCarType + ` ORDER BY name`)
	if err != nil {
		_ = db.Close()
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var carTypeId int
		var carTypeName string
		err = rows.Scan(&carTypeId, &carTypeName)
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, entities.CarType{
			Name:              carTypeName,
			Id:                carTypeId,
		})
	}
	err = rows.Err()
	if err != nil {
		_ = db.Close()
		panic(err)
	}
	_ = db.Close()
	return nodes
}

func FindAllBookings(typeId int) []entities.CarBooking{
	var bookings = []entities.CarBooking{}

	db := getDatabaseClient()

	var rows *sql.Rows
	var err error

	if typeId != -1 {
		rows, err = db.Query(`SELECT cb.id, cb.car_id, cct.id, cct.name, cb.supplier, cb.time, cb.departure_id, cn.name, cb.arrival_id, cn2.name FROM ` + CollectionBooking +` cb 
									INNER JOIN ` + CollectionCar + ` cc ON (cc.id = cb.car_id) 
									INNER JOIN ` + CollectionCarType + ` cct ON (cct.id = cc.car_type)
									INNER JOIN ` + CollectionNode + ` cn ON (cn.id = cb.departure_id)
									INNER JOIN ` + CollectionNode + ` cn2 ON (cn2.id = cb.arrival_id)
									WHERE cct.id = $1`, typeId)
	}else{
		rows, err = db.Query(`SELECT cb.id, cb.car_id, cct.id, cct.name, cb.supplier, cb.time, cb.departure_id, cn.name, cb.arrival_id, cn2.name FROM ` + CollectionBooking +` cb 
									INNER JOIN ` + CollectionCar + ` cc ON (cc.id = cb.car_id) 
									INNER JOIN ` + CollectionCarType + ` cct ON (cct.id = cc.car_type)
									INNER JOIN ` + CollectionNode + ` cn ON (cn.id = cb.departure_id)
									INNER JOIN ` + CollectionNode + ` cn2 ON (cn2.id = cb.arrival_id)`)
	}

	if err != nil {
		_ = db.Close()
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var bookId int
		var carId int
		var typeId int
		var typeName string
		var supplier string
		var timestamp int
		var departureId int
		var departureName string
		var arrivalId int
		var arrivalName string
		err = rows.Scan(&bookId, &carId, &typeId, &typeName, &supplier, &timestamp, &departureId, &departureName, &arrivalId, &arrivalName)
		if err != nil {
			panic(err)
		}
		bookings = append(bookings, entities.CarBooking{
			Supplier:  supplier,
			Date:      time.Time{},
			Id:        bookId,
			Arrival:   entities.Node{
				Name:              arrivalName,
				Id:                arrivalId,
				AvailableCarTypes: GetCarTypesForNode(arrivalId),
			},
			Departure: entities.Node{
				Name:              departureName,
				Id:                departureId,
				AvailableCarTypes: GetCarTypesForNode(departureId),
			},
			Car:       entities.Car{
				Id:      carId,
				CarType: entities.CarType{
					Name: typeName,
					Id:   typeId,
				},
			},
		})
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		_ = db.Close()
		panic(err)
	}
	_ = db.Close()
	return bookings
}

func GetNodeFromId(id int) (entities.Node, error){
	db := getDatabaseClient()
	row := db.QueryRow(`SELECT cn.id, cn.name FROM ` + CollectionNode +` cn WHERE cn.id = $1`, id)

	var nodeId int
	var nodeName string

	switch err := row.Scan(&nodeId, &nodeName); err {
	case sql.ErrNoRows:
		_ = db.Close()
		return entities.Node{}, errors.New("No node for id " + strconv.Itoa(id))
	case nil:
		break
	default:
		panic(err)
	}

	_ = db.Close()

	return entities.Node{
		Id: nodeId,
		Name: nodeName,
		AvailableCarTypes: GetCarTypesForNode(nodeId),
	}, nil
}

func GetCarTypesForNode(nodeId int) []entities.CarType{
	db := getDatabaseClient()
	rows, err := db.Query(`SELECT cnt.car_type_id, ct.name FROM ` + CollectionNodeCarType +` cnt INNER JOIN ` + CollectionCarType + ` ct ON (ct.id = cnt.car_type_id) WHERE cnt.node_id = $1 ORDER BY cnt.car_type_id`, nodeId)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	var types []entities.CarType
	defer rows.Close()
	for rows.Next() {
		var typeId int
		var typeName string
		err = rows.Scan(&typeId, &typeName)
		if err != nil {
			_ = db.Close()
			panic(err)
		}
		types = append(types, entities.CarType{
			Name: typeName,
			Id:   typeId,
		})
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		_ = db.Close()
		panic(err)
	}
	_ = db.Close()
	return types
}

func GetCarFromId(id int) (entities.Car, error) {
	db := getDatabaseClient()
	row := db.QueryRow(`SELECT c.id, ct.id, ct.name FROM ` + CollectionCar +` c INNER JOIN ` + CollectionCarType + ` ct ON (c.car_type = ct.id) WHERE c.id = $1`, id)
	var carId int
	var typeId int
	var typeName string

	switch err := row.Scan(&carId, &typeId, &typeName); err {
		case sql.ErrNoRows:
			_ = db.Close()
			return entities.Car{}, errors.New("No car for id " + strconv.Itoa(id))
		case nil:
			break
		default:
			panic(err)
	}

	_ = db.Close()
	return entities.Car{
		Id: carId,
		CarType: entities.CarType{
			Name: typeName,
			Id:   typeId,
		},
	}, nil
}
