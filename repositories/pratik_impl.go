package repositories

import (
	"context"
	"errors"
	"first/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
    collection := client.Database("pratik").Collection("pratik")
    fmt.Println("COLLECTION PRATIK INITIALIZED")

    // Perform a no-op write to ensure the collection is created
    _, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Options: options.Index().SetName("noop"),
    })
    if err != nil {
        fmt.Println("Error creating no-op index:", err)
    }

    return &UserRepository{
        collection: collection,
    }
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (primitive.ObjectID, error) {
    result, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return primitive.NilObjectID, err
    }
    insertedId, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        return primitive.NilObjectID, errors.New("inserted ID is not an ObjectId")
    }
    return insertedId, nil
}
