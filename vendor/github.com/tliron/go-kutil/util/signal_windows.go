//go:build windows

package util

import (
	"os"
)

var shutdownSignals = []os.Signal{os.Interrupt}
