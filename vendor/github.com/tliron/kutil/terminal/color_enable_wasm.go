//go:build wasm

package terminal

import (
	"os"
)

// Returns false.
func EnableColor(file *os.File) (bool, CleanupFunc, error) {
	return false, nil, nil
}
