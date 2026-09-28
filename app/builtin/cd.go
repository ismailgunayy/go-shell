package builtin

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"syscall"
)

var lastDir string

func Cd(streams Streams, params []string) (exitCode int, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return 1, err
	}

	if len(params) == 0 {
		return 1, errors.New("cd: needs a parameter")
	}

	if len(params) > 1 {
		return 1, errors.New("cd: too many parameters")
	}

	dir := params[0]

	if dir == "-" && lastDir != "" {
		dir = lastDir
	}

	if strings.HasPrefix(dir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return 1, err
		}

		dir = strings.Replace(dir, "~", homeDir, 1)
	}

	if err := os.Chdir(dir); err != nil {
		if _, ok := errors.AsType[*fs.PathError](err); ok {
			switch {
			case os.IsNotExist(err):
				return 1, fmt.Errorf("cd: %s: No such file or directory", dir)
			case os.IsPermission(err):
				return 1, fmt.Errorf("cd: %s: Permission denied", dir)
			case errors.Is(err, syscall.ENOTDIR):
				return 1, fmt.Errorf("cd: %s: Not a directory", dir)
			}
		}
	} else {
		lastDir = currentDir
	}

	return 0, nil
}
