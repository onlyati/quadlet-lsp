//go:build !windows

package util

import (
	"os"
	"syscall"
)

func ExitOnSignals() {
	ExitOnSignal(os.Interrupt, 130)    // CTRL+C
	ExitOnSignal(syscall.SIGTERM, 143) // "kill"
}
