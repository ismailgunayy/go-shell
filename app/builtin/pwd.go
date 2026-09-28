package builtin

import (
	"fmt"
	"os"
)

func Pwd(streams Streams, params []string) (exitCode int, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return 1, err
	}

	fmt.Fprintln(streams.Out, dir)
	return 0, nil
}
