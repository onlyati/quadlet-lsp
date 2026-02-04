//go:build !wasm

package terminal

import (
	"os"

	"golang.org/x/term"
)

func GetSize() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}
