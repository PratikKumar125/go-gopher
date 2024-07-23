package repositories

import (
	"first/repositories/user_repository"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
    PratikRepo *user_repository.UserRepository
}

func RegisterRepositories(container *dig.Container) error {
    if err := container.Provide(NewDBClient); err != nil {
        return err
    }
    if err := container.Provide(func(dbClient *MongoClient) *user_repository.UserRepository {
        return user_repository.NewUserRepository(dbClient.client)
    }); err != nil {
        return err
    }
    return nil
}
