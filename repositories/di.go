package repositories

import (
	"first/utils"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
    PratikRepo *UserRepository
}

func RegisterRepositories(container *dig.Container) error {
    if err := container.Provide(func(client *utils.MongoClient) *UserRepository {
        return NewUserRepository(client.Client())
    }); err != nil {
        return err
    }
    return nil
}
