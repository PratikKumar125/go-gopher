package utils

import (
	"first/tasks"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
    dig.In
    Cache      *Cache
    MongoClient *MongoClient
    AsynqClientStruct *AsynqClient
    AsynqServerStruct *AsynqServer
    TaskHandlerStruct *tasks.HandlerStruct
}

func RegisterDependencies(container *dig.Container) error {
    if err := container.Provide(NewCache); err != nil {
        return err
    }
    if err := container.Provide(NewClient); err != nil {
        return err
    }
    if err := container.Provide(NewAsynqClient); err != nil {
        return err
    }
    if err := container.Provide(NewAsynqServer); err != nil {
        return err
    }
    if err := container.Provide(func(client *MongoClient) *tasks.HandlerStruct {
        return tasks.NewHandler(client.Client())
    }); err != nil {
        return err
    }
    return nil
}
