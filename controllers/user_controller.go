package controllers

import (
	"errors"
	"first/models"
	"first/repositories"
	"first/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (ctrl *UserController) ProtectedUser(ctx *fiber.Ctx) error {
	//Handling user from request object, coming from middleware SuccessHandler, this is just for the demo
	token_user := ctx.Locals("user").(*jwt.Token)
	jwt_claim := token_user.Claims.(jwt.MapClaims)
	name := jwt_claim["name"].(string)
	fmt.Println("Name of user from decoded token is:", name)
	return ctx.JSON(fiber.Map{
		"hello": name,
	})
}