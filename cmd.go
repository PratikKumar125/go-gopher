package main

import (
	"first/commands"
	"fmt"
)

func execute() {
	all_commands := commands.GetAllCommands()
	fmt.Println(all_commands, "ALL AVAILABLE COMMANDS")
	st, err := commands.GetCommandExecutableFn("run-dummy-command")
	if err != nil {
		panic(err)
	}
	st()
	fmt.Println("COMMAND RUN DONE")
}
