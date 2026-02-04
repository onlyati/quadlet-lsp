package util

import (
	"os"

	"golang.org/x/term"
)

func ToRawTerminal(f func() error) {
	stdin := int(os.Stdin.Fd())
	state, err := term.MakeRaw(stdin)
	FailOnError(err)

	handle := OnExitError(func() error {
		os.Stdout.WriteString("\r")
		return term.Restore(stdin, state)
	})

	err = f()

	if err == nil {
		Exit(0)
	} else {
		handle.Cancel()
		os.Stdout.WriteString("\r")
		term.Restore(stdin, state)
		FailOnError(err)
	}
}
