package util

import (
	"fmt"

	"github.com/tliron/kutil/terminal"
)

func Fail(message string) {
	if !terminal.Quiet {
		terminal.Eprintln(terminal.StderrStylist.Error(message))
	}
	Exit(1)
}

func Failf(f string, args ...any) {
	Fail(fmt.Sprintf(f, args...))
}

func FailOnError(err error) {
	if err != nil {
		Fail(err.Error())
	}
}
