# go-shell

A small POSIX-ish shell in Go.

```sh
go run ./app
```

## Features

- Builtins: `echo`, `pwd`, `cd` (`~`, `-`), `type`, `exit`, `complete`
- Runs executables found on `$PATH`
- Single quotes, double quotes, backslash escapes
- Redirects: `>`, `>>`, `2>`, `2>>`
- `$VAR` and `$?` expansion
- Tab completion for commands, file paths, and `complete -C` scripts

## Limitations

- No pipes, `&&`/`||`, `;`, globbing, subshells or job control
- `$VAR` expands only when it is a whole word (`$HOME/x` doesn't work), and it also expands inside single quotes
- `exit` ignores its status code
- `complete -C` scripts get no arguments or `COMP_*` variables

## Known issues

- A quoted `'>'` is still treated as a redirect
- An external command that fails to start fails silently (`$?` = -1)
