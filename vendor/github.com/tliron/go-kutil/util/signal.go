package util

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

// Registers handlers for SIGINT and (on Posix systems) SIGTERM.
// The returned channel will be closed when either signal is sent.
func SetupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	stopChannel := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stopChannel)
		<-c
		Exit(1) // second signal. Exit directly.
	}()

	return stopChannel
}
