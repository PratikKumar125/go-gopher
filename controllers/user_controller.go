package controllers

import (
	"errors"
	"first/models"
	"first/repositories"
	"first/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userRepo *repositories.UserRepository
	Cache *utils.Cache
}

func NewUserController(user_repo *repositories.UserRepository, cache *utils.Cache) *UserController {
	fmt.Println("User controller initialized")
	return &UserController{
		userRepo: user_repo,
		Cache: cache,
	}
}

func (ctrl *UserController) CreateNewUser(ctx  *fiber.Ctx) (error) {
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
	oid, err := ctrl.userRepo.CreateUser(ctx.Context(), user); if err != nil {
		return errors.New("failed to create user, try again")
	}
	fmt.Println("User created with ObjectId as", oid)
	ctrl.Cache.Client().Set(ctx.Context(), "ping", "pong", 0)
	fmt.Println("SET TO CACHE DONE")
	return ctx.JSON(fiber.Map{
		"id": oid.Hex(),
	})
}