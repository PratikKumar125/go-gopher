package main

import (
	"first/commands"
	"first/di"
	"fmt"
)

func main() {
	if err := di.InitCommands(); err != nil {
		fmt.Println("Failed to initialize dependencies:", err)
		return
	}

	err := di.CommandsContainer.Invoke(func(inj *di.CommandsInjected) {
		func () {
			inj.Commands.DummyCommandStruct.RegisterDummyCommand()
		} ()
	})
	if err != nil {
		panic((err))
	}


	all_commands := commands.GetAllCommands()
	fmt.Println(all_commands, "ALL AVAILABLE COMMANDS")
	st, err := commands.GetCommandExecutableFn("run-dummy-command")
	if err != nil {
		panic(err)
	}
	st()
	fmt.Println("COMMAND RUN DONE")
}
