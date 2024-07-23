package routes

import (
	"first/controllers/users"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	Router *Router
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(func (UserController *users.UserController) *Router {
		return NewRouter(UserController)
	}); err != nil {
		return err
	}
	return nil
}
