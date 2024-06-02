package tasks

import (
	"context"
	"encoding/json"
	"first/models"
	"fmt"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerStruct struct {
	db *mongo.Client
}

func NewHandler(client *mongo.Client) *HandlerStruct {
	return &HandlerStruct{
		db: client,
	}
}

func (hs *HandlerStruct) HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil { // Access t.Payload directly
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	var user models.User
	err := hs.db.Database("pratik").Collection("pratik").FindOne(context.Background(), bson.M{"_id": p.UserID}).Decode(&user)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %v", err)
	}

	fmt.Println("EMAIL SENT TO USER", user)
	return nil
}
