package executor

import (
	"fmt"
	"go-shell/app/tokenizer"
	"os"
	"slices"
)

type Redirect struct {
	index int
	kind  tokenizer.RedirectKind
}

func (e *Executor) handleRedirect() (exitCode int, err error) {
	redirects := []Redirect{}

	for i, v := range e.Params {
		var redirect *Redirect

		if slices.Contains(tokenizer.RedirectOperators, tokenizer.RedirectKind(v)) {
			redirect = &Redirect{
				index: i,
				kind:  tokenizer.RedirectKind(v),
			}
		}

		if redirect != nil {
			if i+1 >= len(e.Params) || slices.Contains(tokenizer.RedirectOperators, tokenizer.RedirectKind(e.Params[i+1])) {
				return 1, fmt.Errorf("syntax error near unexpected token %s", e.Params[redirect.index])
			}

			redirects = append(redirects, *redirect)
		}
	}

	if len(redirects) > 0 {
		for i, redirect := range redirects {
			// Since we are deleting items from params during the loop, we need to readjust the index of redirects
			redirect.index -= i * 2
			fileName := e.Params[redirect.index+1]
			e.Params = append(e.Params[:redirect.index], e.Params[redirect.index+2:]...)
			var fileFlag int

			switch redirect.kind {
			case tokenizer.RedirectOut, tokenizer.RedirectErr:
				fileFlag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC

			case tokenizer.RedirectAppendOut, tokenizer.RedirectAppendErr:
				fileFlag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
			}

			file, err := os.OpenFile(fileName, fileFlag, 0644)
			if err != nil {
				return 1, err
			}

			e.files = append(e.files, file)

			switch redirect.kind {
			case tokenizer.RedirectOut, tokenizer.RedirectAppendOut:
				e.Streams.Out = file

			case tokenizer.RedirectErr, tokenizer.RedirectAppendErr:
				e.Streams.Err = file
			}
		}
	}

	return 0, nil
}
