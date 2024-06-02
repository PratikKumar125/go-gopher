package commands

type MyFuncType func()

var (
	cmds = []string{}
	m = make(map[string]MyFuncType)
)

func GetAllCommands() []string {
	return cmds
}

func GetCommandExecutableFn(cmd string) (MyFuncType, error) {
	_, ok := m[cmd]
	if ok {
		return m[cmd], nil
	} 
	return nil, nil
}

func CreateNewCommand(cmd string, execute MyFuncType) error {
	cmds = append(cmds, cmd)
	m[cmd] = execute
	return nil
}
