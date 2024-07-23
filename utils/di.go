package utils

import (
	"go.uber.org/dig"
)

type DependenciesHolder struct {
    dig.In
    Cache      *Cache
    AsynqClientStruct *AsynqClient
    AsynqServerStruct *AsynqServer
}

func RegisterDependencies(container *dig.Container) error {
    if err := container.Provide(NewCache); err != nil {
        return err
    }
    if err := container.Provide(NewAsynqClient); err != nil {
        return err
    }
    if err := container.Provide(NewAsynqServer); err != nil {
        return err
    }
    return nil
}
