package util

import (
	"fmt"
	"os"
	signalpkg "os/signal"
)

// Inspired by: https://github.com/tebeka/atexit

var exitHooks []exitHook
var exitNextHandle ExitFunctionHandle = ExitFunctionHandle(1)
var exitLock RWLocker = NewSyncRWLocker()

type exitHook struct {
	function func()
	handle   ExitFunctionHandle
}

func OnExit(exitFunction func()) ExitFunctionHandle {
	exitLock.Lock()
	defer exitLock.Unlock()

	handle := exitNextHandle
	exitNextHandle++

	exitHooks = append(exitHooks, exitHook{
		function: exitFunction,
		handle:   handle,
	})

	return handle
}

func OnExitError(exitFunction func() error) ExitFunctionHandle {
	return OnExit(func() {
		if err := exitFunction(); err != nil {
			fmt.Fprintf(os.Stderr, "error during exit: %s\n", err.Error())
		}
	})
}

func Exit(code int) {
	exitLock.RLock()

	length := len(exitHooks)
	for index := range exitHooks {
		exitHook := exitHooks[length-index-1] // reverse order
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Fprintf(os.Stderr, "panic during exit: %+v\n", r)
				}
			}()

			exitHook.function()
		}()
	}

	exitLock.RUnlock()

	os.Exit(code)
}

func ExitOnSignal(signal os.Signal, exitCode int) {
	ch := make(chan os.Signal)
	signalpkg.Notify(ch, signal) // CTRL+C
	go func() {
		<-ch
		Exit(exitCode)
	}()
}

//
// ExitFunctionHandle
//

type ExitFunctionHandle int

func (self ExitFunctionHandle) Cancel() {
	exitLock.Lock()
	defer exitLock.Unlock()

	for index, exitEntry := range exitHooks {
		if exitEntry.handle == self {
			exitHooks = append(exitHooks[:index], exitHooks[index+1:]...)
			break
		}
	}
}
