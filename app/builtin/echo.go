package builtin

import (
	"fmt"
	"strings"
)

func Echo(streams Streams, params []string) (exitCode int, err error) {
	fmt.Fprintln(streams.Out, strings.Join(params, " "))
	return 0, nil
}
