//go:build !windows && !wasm

package terminal

import (
	"os"

	"github.com/muesli/termenv"
)

// Returns true if the file supports colorization.
func EnableColor(file *os.File) (bool, CleanupFunc, error) {
	output := termenv.NewOutput(file)
	colorize := output.EnvColorProfile() != termenv.Ascii
	return colorize, nil, nil
}
