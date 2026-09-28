package completer

import (
	"go-shell/app/builtin"
	"go-shell/app/completer/trie"
	"os"
	"runtime"
	"strings"
	"sync"
)

func (c *Completer) initCommandsAndExecutables() {
	c.availableCommands = trie.NewTrie()

	// Builtin commands
	for k := range builtin.ShellCommands {
		c.availableCommands.Insert([]rune(k))
	}

	// PATH executables
	path := os.Getenv("PATH")
	pathDirs := strings.SplitSeq(path, string(os.PathListSeparator))

	var wg sync.WaitGroup
	tasks := make(chan string)
	results := make(chan string)

	for range runtime.NumCPU() {
		wg.Go(func() {
			for dir := range tasks {
				files, err := os.ReadDir(dir)
				if err != nil {
					continue
				}

				for _, file := range files {
					fileInfo, err := file.Info()
					if err != nil {
						continue
					}

					isExecutable := !fileInfo.IsDir() && fileInfo.Mode().Perm()&0o111 != 0
					if isExecutable {
						results <- fileInfo.Name()
					}
				}
			}
		})
	}

	go func() {
		for dir := range pathDirs {
			tasks <- dir
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for executable := range results {
		c.availableCommands.Insert([]rune(executable))
	}
}
