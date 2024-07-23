package user_repository

import (
	"context"
	"errors"
	"first/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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

func (r * UserRepository) FindOneUser(ctx context.Context, email string) (models.User, error) {
    filter := bson.D{{"email", email}}
    userStruct := new(models.User)
    err := r.collection.FindOne(ctx, filter).Decode(&userStruct)
    if err != nil {
        if err == mongo.ErrNoDocuments {
		    return *userStruct, err
	    }
	    return *userStruct, fmt.Errorf("failed to find users: %w", err)
    }
    return *userStruct, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
    findOptions := options.Find()
    findOptions.SetLimit(5)
    var users []models.User
    cur, err := r.collection.Find(ctx, bson.D{{}}, findOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to find users: %w", err)
    }
    defer cur.Close(ctx)

    for cur.Next(ctx) {
        var elem models.User
        if err := cur.Decode(&elem); err != nil {
            return nil, fmt.Errorf("failed to decode user: %w", err)
        }
        users = append(users, elem)
    }

    if err := cur.Err(); err != nil {
        return nil, fmt.Errorf("cursor error: %w", err)
    }

    return users, nil
}
