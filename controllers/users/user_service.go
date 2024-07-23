package users

import (
	"errors"
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
	dependencies.CacheClient.Client().Set(ctx.Context(), "ping", "pong", 0)
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
	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (dependencies *UserServiceStruct) HandleGetUserProfile(ctx *fiber.Ctx) (error) {
	//Handling user from request object, coming from middleware SuccessHandler, 
	//this is just for the demo
	token_user := ctx.Locals("user").(*jwt.Token)
	jwt_claim := token_user.Claims.(jwt.MapClaims)
	name := jwt_claim["name"].(string)
	fmt.Println("Name of user from decoded token is:", name)
	user, err := dependencies.UserRepo.FindOneUser(ctx.Context(), jwt_claim["email"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return ctx.JSON(fiber.Map{
		"user": user,
	})
}