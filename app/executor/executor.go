package executor

import (
	"errors"
	"fmt"
	"go-shell/app/builtin"
	"os"
	"os/exec"
)

type Executor struct {
	Command string
	Params  []string
	Streams builtin.Streams
	State   *State
	files   []*os.File
}

func (e *Executor) Execute() {
	exitCode, err := e.executeWithError()
	if err != nil {
		fmt.Fprintf(e.Streams.Err, "%s\n", err.Error())
	}

	e.State.LastExitCode = exitCode
}

func (e *Executor) executeWithError() (exitCode int, err error) {
	exitCode, err = e.handleRedirect()
	if err != nil {
		return exitCode, err

	}

	defer func() {
		for _, f := range e.files {
			f.Close()
		}
	}()

	op, builtIn := builtin.ShellCommands[builtin.CommandName(e.Command)]

	if builtIn {
		code, err := op(e.Streams, e.Params)
		if err != nil {
			return code, err
		}
	} else {
		_, err := exec.LookPath(e.Command)
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				return 1, fmt.Errorf("%s: command not found", e.Command)
			} else {
				return 1, err
			}
		}

		cmd := exec.Command(e.Command, e.Params...)
		cmd.Stdin = e.Streams.In
		cmd.Stdout = e.Streams.Out
		cmd.Stderr = e.Streams.Err
		if err := cmd.Run(); err != nil {
			if _, ok := errors.AsType[*exec.ExitError](err); !ok {
				return 1, err
			}
		}

		return cmd.ProcessState.ExitCode(), nil
	}

	return 0, nil
}
