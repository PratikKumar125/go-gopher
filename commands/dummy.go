package commands

import (
	"context"
	"first/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DummyStruct struct {
	DB *mongo.Client
}

func NewDummyCommand(db *mongo.Client) *DummyStruct {
	return &DummyStruct{
		DB: db,
	}
}

func (ds *DummyStruct) Execute() {
	var operations []mongo.WriteModel

	admin1 := models.User{
		Name: "ADMIN1",
		Email: "admin1@admin.com",
	}

	OpeartionAdmin1 := mongo.NewUpdateOneModel()
	OpeartionAdmin1.SetFilter(bson.M{"email": admin1.Email})
	OpeartionAdmin1.SetUpdate(bson.M{"$set": bson.M{"name": admin1.Name, "email": admin1.Email}})
	OpeartionAdmin1.SetUpsert(true)
	operations = append(operations, OpeartionAdmin1)

	coll := ds.DB.Database("pratik").Collection("pratik")
	opts := options.BulkWrite().SetOrdered(true)
	result, err := coll.BulkWrite(context.TODO(), operations, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DUMMY COMMAND WORKED SUCCESFULLY")
	fmt.Println("ADMIN SEEDED INTO DATABASE", result)
}

func (ds *DummyStruct) RegisterDummyCommand() {
	err := CreateNewCommand("run-dummy-command", ds.Execute)
	if err != nil {
		fmt.Println("Failed to create new command:", err)
	}
	fmt.Println("DUMMY COMMAND REGISTERED")
}