package di

import (
	"first/commands"

	"go.uber.org/dig"
)

type CommandsInjected struct {
    Commands commands.DependenciesHolder
}

func NewCommandsInjected(
    cm commands.DependenciesHolder,
) *CommandsInjected {
    return &CommandsInjected{
        Commands: cm,
    }
}

var CommandsContainer *dig.Container

func InitCommands() error {
	CommandsContainer = dig.New()
  if err := commands.RegisterConsoleCommands(CommandsContainer); err != nil {
    return err
  }
  if err := CommandsContainer.Provide(NewCommandsInjected); err != nil {
    return err
  }
  return nil
}
