package tokenizer

import (
	"slices"
	"strings"
)

type mode int

const (
	modeNormal mode = iota
	modeSingleQuote
	modeDoubleQuote
)

const (
	SingleQuote = '\''
	DoubleQuote = '"'
	Space       = ' '
	BackTick    = '`'
	NewLine     = '\n'
	Variable    = '$'
	BackSlash   = '\\'
	redirect    = '>'
)

type RedirectKind string

const (
	RedirectOut       RedirectKind = "1>"
	RedirectErr       RedirectKind = "2>"
	RedirectAppendOut RedirectKind = "1>>"
	RedirectAppendErr RedirectKind = "2>>"
)

var RedirectOperators = []RedirectKind{RedirectOut, RedirectErr, RedirectAppendOut, RedirectAppendErr}

var escapedChars = []rune{DoubleQuote, BackSlash, Variable, BackTick, NewLine}

type tokenizer struct {
	mode    mode
	escaped bool
	token   strings.Builder
	tokens  []string
	i       int
}

func Tokenize(input []rune) (tokens []string) {
	t := tokenizer{mode: modeNormal}

	for ; t.i < len(input); t.i++ {
		switch t.mode {
		case modeSingleQuote:
			t.stepSingleQuote(input)

		case modeDoubleQuote:
			t.stepDoubleQuote(input)

		default:
			t.stepNormal(input)
		}
	}
	t.flush()

	return t.tokens
}

func (t *tokenizer) stepSingleQuote(chars []rune) {
	if chars[t.i] == SingleQuote {
		t.mode = modeNormal
	} else {
		t.token.WriteRune(chars[t.i])
	}
}

func (t *tokenizer) stepDoubleQuote(chars []rune) {
	switch {
	case t.escaped:
		if !slices.Contains(escapedChars, chars[t.i]) {
			t.token.WriteRune('\\')
		}

		t.token.WriteRune(chars[t.i])
		t.escaped = false

	case chars[t.i] == BackSlash:
		t.escaped = true

	case chars[t.i] == DoubleQuote:
		t.mode = modeNormal

	default:
		t.token.WriteRune(chars[t.i])
	}
}

func (t *tokenizer) stepNormal(chars []rune) {
	if t.escaped {
		t.token.WriteRune(chars[t.i])
		t.escaped = false
		return
	}

	switch chars[t.i] {
	case SingleQuote:
		t.mode = modeSingleQuote

	case DoubleQuote:
		t.mode = modeDoubleQuote

	case BackSlash:
		t.escaped = true

	case redirect:
		if tStr := t.token.String(); len(tStr) > 0 && tStr != "1" && tStr != "2" {
			t.flush()
		}

		tStr := t.token.String()
		if tStr == "" {
			tStr = "1"
		}

		op := tStr + ">"
		if t.i+1 < len(chars) && chars[t.i+1] == '>' {
			op += ">"
			t.i++
		}

		t.token.Reset()
		t.token.WriteString(op)
		t.flush()

	case Space:
		t.flush()

	default:
		t.token.WriteRune(chars[t.i])
	}
}

func (t *tokenizer) flush() {
	if t.token.Len() > 0 {
		t.tokens = append(t.tokens, t.token.String())
		t.token.Reset()
	}
}
