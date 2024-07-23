package commands

import (
	"first/repositories"
	"fmt"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	DummyCommandStruct *DummyStruct
    // DbClient *repositories.MongoClient
}

func RegisterConsoleCommands(container *dig.Container) error {
    if err := container.Provide(repositories.NewDBClient()); err != nil {
        return err
    }
    if err := container.Provide(func() *DummyStruct {
        client := repositories.NewDBClient().Client()
        return NewDummyCommand(client)
    }); err != nil {
        fmt.Println("FAILED", repositories.NewDBClient().Client())
        return err
    }
    return nil
}
