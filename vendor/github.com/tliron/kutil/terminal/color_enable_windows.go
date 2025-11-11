//go:build windows

package terminal

import (
	"os"

	"github.com/muesli/termenv"
)

// Checks if the file supports colorization, and if it does then will
// attempt to enable ANSI on it, returning an error if it fails. If it
// succeeds will return a [CleanupFunc] that restores the file to its
// original state.
func EnableColor(file *os.File) (bool, CleanupFunc, error) {
	output := termenv.NewOutput(file)
	if output.EnvColorProfile() != termenv.Ascii {
		if mode, err := termenv.EnableWindowsANSIConsole(); err == nil {
			return true, func() error {
				return termenv.RestoreWindowsConsole(mode)
			}, nil
		} else {
			return false, nil, err
		}
	} else {
		return false, nil, nil
	}
}
