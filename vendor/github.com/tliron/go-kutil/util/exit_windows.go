//go:build windows

package util

import (
	"os"
)

func ExitOnSignals() {
	ExitOnSignal(os.Interrupt, 130) // CTRL+C
}
