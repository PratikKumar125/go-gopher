package users

import (
	"errors"
	"first/controllers/exceptions"
	"first/controllers/transformers"
	"first/models"
	"first/repositories/user_repository"
	"first/tasks"
	"first/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
)

type UserServiceStruct struct {
	UserRepo *user_repository.UserRepository
	CacheClient *utils.Cache
	QueueClient *utils.AsynqClient
}

func NewServiceStruct(userRepo *user_repository.UserRepository, cache *utils.Cache, queueClient *utils.AsynqClient) *UserServiceStruct{
	return &UserServiceStruct{
		UserRepo: userRepo,
		CacheClient: cache,
		QueueClient: queueClient,
	}
}

func (dependencies *UserServiceStruct) HandleCreateNewUser(ctx *fiber.Ctx) (error) {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(user); err != nil {
		return utils.CheckForValidation(ctx, err, fiber.StatusUnprocessableEntity, "user")
	}

	oid, err := dependencies.UserRepo.CreateUser(ctx.Context(), user); if err != nil {
		return errors.New("failed to create user, try again")
	}
	fmt.Println("User created with ObjectId as", oid)
	dependencies.CacheClient.AddKeyWithTTL(ctx.Context(), "ping", "pong", "30s")
	fmt.Println("SET TO CACHE DONE")

	//signing JWT token
	claims := jwt.MapClaims{
		"_id": oid,
		"name": user.Name,
		"email": user.Email,
	}
	token, err := utils.SignJwtToken(&claims) 
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	task1, err := tasks.NewWelcomeEmailTask(oid)
	if err != nil {
		return err
	}
	if _, err := dependencies.QueueClient.Client().Enqueue(
		task1,
		asynq.Queue("critical"),
	); err != nil {
		log.Fatal(err)
	}
	return transformers.GlobalSuccessResponse(ctx, TranformCreateUser(ctx.Context(), token))
}

func (dependencies *UserServiceStruct) HandleGetUserProfile(ctx *fiber.Ctx) error {
	tokenUserInterface := ctx.Locals("user")
	tokenUser := tokenUserInterface.(models.User)
	
	email := tokenUser.Email

	user, err := dependencies.UserRepo.FindOneUser(ctx.Context(), email)
	if err != nil {
		return exceptions.ThrowInternalServerError(ctx)
	}
	return transformers.GlobalSuccessResponse(ctx, TransformUser(ctx.Context(), user))
}

func (dependencies *UserServiceStruct) HandleGetAllUserPaginated(ctx *fiber.Ctx) error {
	user, err := dependencies.UserRepo.FindAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	return transformers.GlobalSuccessResponse(ctx, user)
}
