package builtin

import (
	"fmt"
	"strings"
)

const (
	flagPrintSpec    = "p"
	flagRemoveSpec   = "r"
	flagCompleteSpec = "C"
)

var Specifications = make(map[string]string)

// -p <command-name> | Print completion specification
// -C <path> <command-name> | Add completion specification
// -r <command-name> | Remove completion specification
func Complete(streams Streams, params []string) (exitCode int, err error) {
	// Currently it's assumed that the flag and its argument will always be the first 2 parameters
	if len(params) == 0 {
		return 1, fmt.Errorf("complete: <flag> <path-if-any> <command-name>")
	}

	flagRaw := params[0]
	flagRaw, ok := strings.CutPrefix(flagRaw, "-")
	if !ok {
		return 1, fmt.Errorf("complete: couldn't parse flag")
	}

	args := params[1:]

	switch flagRaw {
	case flagPrintSpec:
		if len(args) == 0 {
			return 1, fmt.Errorf("complete: no command specified")
		}

		command := args[0]

		path, ok := Specifications[command]
		if !ok {
			return 1, fmt.Errorf("complete: %s: no completion specification", command)
		}

		fmt.Fprintf(streams.Out, "complete -C '%s' %s\n", path, command)

	case flagCompleteSpec:
		if len(args) < 2 {
			return 1, fmt.Errorf("complete: specify path and command")
		}

		path, command := args[0], args[1]

		if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "./") {
			path = "./" + path
		}

		Specifications[command] = path

	case flagRemoveSpec:
		if len(args) == 0 {
			return 1, fmt.Errorf("complete: no command specified")
		}

		command := args[0]
		delete(Specifications, command)

	default:
		return 1, fmt.Errorf("complete: unknown flag")
	}

	return 0, nil
}
