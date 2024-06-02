package routes

import (
	"first/controllers"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	Router *Router
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(func (UserController *controllers.UserController) *Router {
		return NewRouter(UserController)
	}); err != nil {
		return err
	}
	return nil
}
