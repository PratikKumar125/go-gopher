package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

var (
	URI = "mongodb://localhost:27017/pratik"
)

func NewClient() *MongoClient {

	// Connect to DB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
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