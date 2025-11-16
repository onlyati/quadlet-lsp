// Package commands
//
// This packages contains the callable commands and the command handler
// that can be issued from editors.
package commands

import (
	"errors"
	"sync"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type EditorCommandExecutor struct {
	commands map[string]*allowedCommand
	mutex    sync.Mutex
	rootDir  string
	syncCall bool
}

type allowedCommand struct {
	fn      func(command *protocol.ExecuteCommandParams, e *EditorCommandExecutor, messenger utils.Messenger, executor utils.Commander)
	running bool
}

func NewEditorCommandExecutor(rootDir string) EditorCommandExecutor {
	return EditorCommandExecutor{
		commands: map[string]*allowedCommand{
			"pullAll": {
				fn:      pullAll,
				running: false,
			},
			"listJobs": {
				fn:      listJobs,
				running: false,
			},
			"generateDoc": {
				fn:      GenerateDoc,
				running: false,
			},
		},
		mutex:   sync.Mutex{},
		rootDir: rootDir,
	}
}

func (e *EditorCommandExecutor) Run(command *protocol.ExecuteCommandParams, messenger utils.Messenger, executor utils.Commander) error {
	err := e.tryRun(command, messenger, executor)
	if err != nil {
		messenger.SendMessage(
			utils.MessengerError,
			"Command failed: "+command.Command+", reason: "+err.Error(),
		)
		return nil
	}

	return nil
}

func (e *EditorCommandExecutor) tryRun(command *protocol.ExecuteCommandParams, messenger utils.Messenger, executor utils.Commander) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	v, found := e.commands[command.Command]
	if !found {
		return errors.New("not found")
	}
	if v.running {
		return errors.New("already running")
	}
	v.running = true

	if e.syncCall {
		// Just for unit tests
		e.mutex.Unlock()
		v.fn(command, e, messenger, executor)
		e.mutex.Lock()
		return nil
	}

	go v.fn(command, e, messenger, executor)

	return nil
}

func (e *EditorCommandExecutor) resetRunning(command string) {
	e.mutex.Lock()
	if v, found := e.commands[command]; found {
		v.running = false
	}
	e.mutex.Unlock()
}
