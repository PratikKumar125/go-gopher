package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService *UserServiceStruct
}

func NewUserController(userService *UserServiceStruct) *UserController {
	fmt.Println("User controller initialized")
	return &UserController{
		UserService: userService,
	}
}

func (ctrl *UserController) CreateNewUser(ctx  *fiber.Ctx) (error) {
	return ctrl.UserService.HandleCreateNewUser(ctx)
}

func (ctrl *UserController) ProtectedUser(ctx *fiber.Ctx) error {
	return ctrl.UserService.HandleGetUserProfile(ctx)
}