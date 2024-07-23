package tasks

import (
	"context"
	"encoding/json"
	"first/repositories/user_repository"
	"fmt"

	"github.com/hibiken/asynq"
)

type HandlerStruct struct {
	userRepo *user_repository.UserRepository
}

func NewHandler(userRepo *user_repository.UserRepository) *HandlerStruct {
	return &HandlerStruct{
		userRepo: userRepo,
	}
}

func (hs *HandlerStruct) HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil { // Access t.Payload directly
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	res, err := hs.userRepo.FindOneUser(c, "prateek@gmail.com")
	if err != nil {
		return fmt.Errorf("failed to fetch user: %v", err)
	}

	fmt.Println("EMAIL SENT TO USER", res)
	return nil
}
