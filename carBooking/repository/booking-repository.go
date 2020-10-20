package repository

import (
	"carBooking/entities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"time"
)

var Database *mongo.Database

var CollectionNode = "node"
var CollectionBooking = "booking"
var CollectionCarType = "car_type"
var CollectionCar = "car"

func InitDatabase(){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://" + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT")))
	if err != nil{
		log.Fatal("Error to connect to database")
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	Database = client.Database(os.Getenv("MONGO_DB"))

	collection := Database.Collection(CollectionCarType)
	var carSolid = entities.CarType{Name: "Solid", Id:   1}
	_, _ = collection.InsertOne(ctx, carSolid)
	var carLiquid = entities.CarType{Name: "Liquid", Id:   2}
	_, _ = collection.InsertOne(ctx, carLiquid)

	collection = Database.Collection(CollectionNode)
	var nodeNice = entities.Node{Name:           "Nice", Id:             1, AvailableCarTypes: []entities.CarType{carLiquid}}
	_, _ = collection.InsertOne(ctx, nodeNice)
	var nodeMarseille = entities.Node{Name:           "Marseilles", Id:             2, AvailableCarTypes: []entities.CarType{carLiquid, carSolid}}
	_, _ = collection.InsertOne(ctx, nodeMarseille)

	collection = Database.Collection(CollectionCar)
	var car = entities.Car{Id:             1, CarType: carLiquid}
	_, _ = collection.InsertOne(ctx, car)
	var car2 = entities.Car{Id:             2, CarType: carSolid}
	_, _ = collection.InsertOne(ctx, car2)
}

func CreateBook(Date time.Time, Car entities.Car , Supplier string, NodeDeparture entities.Node, NodeArrival entities.Node){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _ = Database.Collection(CollectionBooking).InsertOne(ctx, entities.CarBooking{
		Date: Date,
		Supplier: Supplier,
		Departure: NodeDeparture,
		Arrival: NodeArrival,
		Car: Car,
	})
}

func FindAllBookings() []entities.CarBooking{
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := Database.Collection(CollectionBooking).Find(ctx, bson.D{})
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)

	var toReturn []entities.CarBooking

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		var carTypesArrival []entities.CarType
		var nodeArrival = result["arrival"].(primitive.M)
		var array = nodeArrival["availablecartypes"].(primitive.A)
		for _, s := range array {
			var theMap = s.(primitive.M)
			carTypesArrival = append(carTypesArrival, entities.CarType{
				Name: theMap["name"].(string),
				Id:   theMap["id"].(int32),
			})
		}

		var carTypeDeparture []entities.CarType
		var nodeDeparture = result["departure"].(primitive.M)
		array = nodeDeparture["availablecartypes"].(primitive.A)
		for _, s := range array {
			var theMap = s.(primitive.M)
			carTypeDeparture = append(carTypeDeparture, entities.CarType{
				Name: theMap["name"].(string),
				Id:   theMap["id"].(int32),
			})
		}

		var car = result["car"].(primitive.M)
		var theMap = car["cartype"].(primitive.M)
		var carType = entities.CarType{
			Name: theMap["name"].(string),
			Id:   theMap["id"].(int32),
		}

		toReturn = append(toReturn, entities.CarBooking{
			Supplier:  result["supplier"].(string),
			Date:      result["date"].(primitive.DateTime).Time(),
			Id:        result["id"].(int32),
			Arrival:   entities.Node{
							Id: nodeArrival["id"].(int32),
							Name: nodeArrival["name"].(string),
							AvailableCarTypes: carTypesArrival,
						},
			Departure: entities.Node{
							Id: nodeDeparture["id"].(int32),
							Name: nodeDeparture["name"].(string),
							AvailableCarTypes: carTypeDeparture,
						},
			Car:       entities.Car{
				Id:      car["id"].(int32),
				CarType: carType,
			},
		})
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return toReturn
}

func GetNodeFromId(id int) (entities.Node, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result bson.M
	err := Database.Collection(CollectionNode).FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return entities.Node{}, errors.New("Error 404: Node with id " + strconv.Itoa(id) + " not found")
	}

	var carTypes []entities.CarType
	var array = result["availablecartypes"].(primitive.A)
	for _, s := range array {
		var theMap = s.(primitive.M)
		carTypes = append(carTypes, entities.CarType{
			Name: theMap["name"].(string),
			Id:   theMap["id"].(int32),
		})
	}

	return entities.Node{
		Id: result["id"].(int32),
		Name: result["name"].(string),
		AvailableCarTypes: carTypes,
	}, nil
}

func GetCarFromId(id int) (entities.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result bson.M
	err := Database.Collection(CollectionCar).FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return entities.Car{}, errors.New("Error 404: Car with id " + strconv.Itoa(id) + " not found")
	}

	var theMap = result["cartype"].(primitive.M)
	var carType = entities.CarType{
		Name: theMap["name"].(string),
		Id:   theMap["id"].(int32),
	}

	return entities.Car{
		Id: result["id"].(int32),
		CarType: carType,
	}, nil
}
