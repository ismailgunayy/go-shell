package completer

import (
	"fmt"
	"go-shell/app/builtin"
	"go-shell/app/completer/trie"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type completionMode int

const (
	commandCompletion completionMode = iota
	specifiedCommandCompletion
	fileNameCompletion
)

type Completer struct {
	state             completerState
	availableCommands *trie.Trie
}

func New() *Completer {
	c := &Completer{}
	c.initCommandsAndExecutables()
	return c
}

func (c *Completer) Do(line []rune, pos int) (completions [][]rune, length int) {
	var longestCommonPrefix []rune
	mode, command, prefix := parseLine(line)
	length = len(prefix)

	switch mode {
	case specifiedCommandCompletion:
		completions, length = handleSpecifiedCommandCompletion(command, prefix)

	case fileNameCompletion:
		completions, length = handleFileNameCompletion(prefix)

	case commandCompletion:
		completions = c.availableCommands.Complete(prefix)
	}

	if slices.Compare(prefix, c.state.lastPrefix) != 0 {
		c.state.lastPrefix = prefix
		c.state.tabPressed = false
	}

	longestCommonPrefix = findLongestCommonPrefix(completions)
	action, newState := decide(len(completions), len(longestCommonPrefix) > 0, c.state)
	c.state = newState

	switch action {
	case actBell:
		completions = nil
		length = 0

		fmt.Fprint(os.Stderr, "\a")

	case actInsertOne:
		only := completions[0]
		if len(only) == 0 || only[len(only)-1] != os.PathSeparator {
			completions[0] = append(completions[0], ' ')
		}

	case actList:
		for i := range completions {
			completions[i] = append(completions[i], ' ')
		}

		names := make([]string, len(completions))
		for i, suffix := range completions {
			if mode != specifiedCommandCompletion {
				names[i] = string(prefix)
			}
			names[i] += strings.TrimRight(string(suffix), " ")
		}

		fmt.Printf("\n%s", strings.Join(names, " "))
		fmt.Printf("\n$ %s", string(line))

	case actInsertLCP:
		completions = [][]rune{longestCommonPrefix}
	}

	return completions, length
}

func parseLine(line []rune) (mode completionMode, command []rune, prefix []rune) {
	lineSplit := strings.Split(strings.TrimLeft(string(line), " "), " ")
	command = []rune(lineSplit[0])

	_, ok := builtin.Specifications[string(command)]
	if ok {
		mode = specifiedCommandCompletion
	} else if len(lineSplit) > 1 {
		mode = fileNameCompletion
	} else {
		mode = commandCompletion
	}

	prefix = []rune(lineSplit[len(lineSplit)-1])
	return mode, command, prefix
}

func handleSpecifiedCommandCompletion(command, prefix []rune) (completions [][]rune, length int) {
	path, ok := builtin.Specifications[string(command)]
	if !ok {
		return
	}

	out, err := exec.Command(path).Output()
	if err != nil {
		return
	}

	out = []byte(strings.Trim(string(out), " \n"))
	lines := strings.SplitSeq(string(out), "\n")

	for line := range lines {
		line, _ = strings.CutPrefix(line, string(prefix))
		if len(line) > 0 {
			completions = append(completions, []rune(line))
		}
	}

	return completions, len(prefix)
}

func handleFileNameCompletion(word []rune) (fileNames [][]rune, length int) {
	isAbsolute := strings.HasPrefix(string(word), string(os.PathSeparator))
	dir := "."

	prefixSplit := strings.Split(string(word), string(os.PathSeparator))
	prefixDir := strings.Join(prefixSplit[:len(prefixSplit)-1], string(os.PathSeparator))
	newWord := prefixSplit[len(prefixSplit)-1]

	if len(prefixDir) > 0 {
		dir = prefixDir
	} else if isAbsolute {
		dir = string(os.PathSeparator)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		fileName := file.Name()

		if file.IsDir() {
			fileName = fileName + string(os.PathSeparator)
		}

		if after, ok := strings.CutPrefix(fileName, string(newWord)); ok {
			fileNames = append(fileNames, []rune(after))
		}
	}

	return fileNames, len(newWord)
}

func findLongestCommonPrefix(completions [][]rune) (prefix []rune) {
	if len(completions) == 0 {
		return
	}

	var longestCommonPrefix strings.Builder
	i := 0

	for {
		var letter rune
		j := 0

		same := true

		for j = range completions {
			if i >= len(completions[j]) {
				same = false
				break
			}

			if letter == 0 {
				letter = completions[j][i]
			}

			if completions[j][i] != letter {
				same = false
				break
			}
		}

		if same {
			longestCommonPrefix.WriteRune(letter)
		} else {
			break
		}

		i++
	}

	return []rune(longestCommonPrefix.String())
}
