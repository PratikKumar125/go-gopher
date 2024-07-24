package user_repository

import (
	"context"
	"errors"
	"first/models"
	"first/repositories/common_repository"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
    collection := client.Database(os.Getenv("MONGO_DATABASE_NAME")).Collection("pratik")
    fmt.Println("COLLECTION USER REPO INITIALIZED")

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

func (r *UserRepository) GetAllPaginated(ctx context.Context, page int32, limit int32) (map[string]interface{}, error) {
    searchPipeline := bson.D{{Key: "$match", Value: bson.D{}}}
    pipeline := common_repository.GetPaginationAggregation(searchPipeline, int(page), int(limit))
    cursor, err := r.collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to execute aggregation: %w", err)
    }
    defer cursor.Close(ctx)

    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
    }

    if len(results) == 0 {
        return map[string]interface{}{
            "meta": map[string]interface{}{
                "total":        0,
                "hasNextPage":  false,
                "hasPrevPage":  false,
            },
            "docs": []models.User{},
        }, nil
    }

    meta := results[0]["meta"].(bson.M)
    docs := results[0]["docs"].(bson.A)

    var users []models.User
    for _, doc := range docs {
        var user models.User
        bsonBytes, _ := bson.Marshal(doc)
        if err := bson.Unmarshal(bsonBytes, &user); err != nil {
            return nil, fmt.Errorf("failed to decode user: %w", err)
        }
        users = append(users, user)
    }

    response := map[string]interface{}{
        "meta": meta,
        "docs": users,
    }

    return response, nil
}
