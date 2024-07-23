package controllers

import (
	"first/controllers/users"
	"first/repositories/user_repository"
	"first/utils"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	UserController *users.UserController
}

func RegisterDependencies(container *dig.Container) error {
	//USERS DEPENDENCIES
	if err := container.Provide(func (userRepo *user_repository.UserRepository, cache *utils.Cache, queueClient *utils.AsynqClient) *users.UserServiceStruct {
		return users.NewServiceStruct(userRepo, cache, queueClient)
	}); err != nil {
		return err
	}
	if err := container.Provide(func (userService *users.UserServiceStruct) *users.UserController {
		return users.NewUserController(userService)
	}); err != nil {
		return err
	}
	//FOLLOW THE SAME FOR OTHER CONTROLLERS
	return nil
}
