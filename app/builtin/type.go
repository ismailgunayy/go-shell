package builtin

import (
	"errors"
	"fmt"
	"os/exec"
)

func Type(streams Streams, params []string) (exitCode int, err error) {
	if len(params) == 0 {
		return 1, errors.New("type: needs a parameter")
	}

	if len(params) > 1 {
		return 1, errors.New("type: many parameters")
	}

	command := params[0]
	_, builtIn := ShellCommands[CommandName(command)]

	if builtIn {
		fmt.Fprintf(streams.Out, "%s is a shell builtin\n", command)
		return 0, nil
	}

	fullPath, err := exec.LookPath(command)
	if err != nil {
		return 1, fmt.Errorf("%s: not found", command)
	} else {
		fmt.Fprintf(streams.Out, "%s is %s\n", command, fullPath)
		return 0, nil
	}
}
