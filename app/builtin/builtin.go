package builtin

import "io"

type CommandName string

const (
	commandExit     CommandName = "exit"
	commandEcho     CommandName = "echo"
	commandType     CommandName = "type"
	commandPwd      CommandName = "pwd"
	commandCd       CommandName = "cd"
	commandComplete CommandName = "complete"
)

type Streams struct {
	In       io.Reader
	Out, Err io.Writer
}

type BuiltinFunc func(streams Streams, params []string) (exitCode int, err error)

var ShellCommands = make(map[CommandName]BuiltinFunc)

func init() {
	ShellCommands = map[CommandName]BuiltinFunc{
		commandExit:     Exit,
		commandEcho:     Echo,
		commandType:     Type,
		commandPwd:      Pwd,
		commandCd:       Cd,
		commandComplete: Complete,
	}
}
