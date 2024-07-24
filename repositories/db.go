package repositories

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

func NewDBClient() *MongoClient {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MONGO DB INITIALIZED")
	return &MongoClient{
		client: client,
	}
}

func (rd *MongoClient) Client() *mongo.Client {
	return rd.client
}