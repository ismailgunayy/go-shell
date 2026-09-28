package builtin

import (
	"os"
)

func Exit(streams Streams, params []string) (exitCode int, err error) {
	os.Exit(0)
	return 0, nil
}
