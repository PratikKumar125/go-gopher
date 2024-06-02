package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TypeWelcomeEmail = "email:welcome"
)

type WelcomeEmailPayload struct {
	UserID primitive.ObjectID `json:"user_id"`
}

func NewWelcomeEmailTask(id primitive.ObjectID) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWelcomeEmail, payload), nil // payload is a []byte
}
