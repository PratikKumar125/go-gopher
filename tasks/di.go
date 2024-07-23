package tasks

import (
	"first/repositories/user_repository"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	Handler *HandlerStruct
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(func (UserRepo *user_repository.UserRepository) *HandlerStruct {
		return NewHandler(UserRepo)
	}); err != nil {
		return err
	}
	return nil
}