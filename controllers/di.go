package controllers

import (
	"first/repositories"
	"first/utils"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	UserController *UserController
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(func (UserRepo *repositories.UserRepository, cache *utils.Cache) *UserController {
		return NewUserController(UserRepo, cache)
	}); err != nil {
		return err
	}
	return nil
}
