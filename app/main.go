package main

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"go-shell/app/builtin"
	"go-shell/app/completer"
	"go-shell/app/executor"
	"go-shell/app/tokenizer"

	"github.com/chzyer/readline"
)

func main() {
	inStream := os.Stdin
	outStream := os.Stdout
	errStream := os.Stderr

	state := &executor.State{}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer.New(),
		EnableMask:   false,
	})
	if err != nil {
		panic(err)
	}

	defer rl.Close()

	for {
		var command string
		var params []string

		input, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			continue
		}
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			panic(err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		tokens := tokenizer.Tokenize([]rune(input))

		// Parameter expansion
		for i := len(tokens) - 1; i >= 0; i-- {
			if after, ok := strings.CutPrefix(tokens[i], "$"); ok {
				if after == "?" {
					tokens[i] = strconv.Itoa(state.LastExitCode)
				} else {
					tokens[i], _ = os.LookupEnv(after)

					if tokens[i] == "" {
						tokens = append(tokens[:i], tokens[i+1:]...)
					}
				}
			}
		}

		if len(tokens) == 0 {
			continue
		}

		command, params = tokens[0], tokens[1:]

		e := (&executor.Executor{
			Command: command,
			Params:  params,
			Streams: builtin.Streams{In: inStream, Out: outStream, Err: errStream},
			State:   state,
		})

		e.Execute()
	}
}
