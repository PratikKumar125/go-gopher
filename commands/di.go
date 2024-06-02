package commands

import (
	"first/utils"
	"fmt"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	DummyCommandStruct *DummyStruct
    DbClient *utils.MongoClient
}

func RegisterConsoleCommands(container *dig.Container) error {
    if err := container.Provide(utils.NewClient); err != nil {
        return err
    }
    if err := container.Provide(func() *DummyStruct {
        client := utils.NewClient().Client()
        return NewDummyCommand(client)
    }); err != nil {
        fmt.Println("FAILED", utils.NewClient().Client())
        return err
    }
    // if err := container.Provide(func() *tasks.HandlerStruct {
    //     client := NewClient().client
    //     return tasks.NewHandler(client)
    // }); err != nil {
    //     return err
    // }
    return nil
}
