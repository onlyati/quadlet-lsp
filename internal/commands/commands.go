package commands

import (
	"errors"
	"sync"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type EditorCommandExecutor struct {
	commands map[string]*allowedCommand
	mutex    sync.Mutex
	rootDir  string
}

type allowedCommand struct {
	fn      func(command string, e *EditorCommandExecutor, ctx glsp.Context, executor utils.Commander)
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
		},
		mutex:   sync.Mutex{},
		rootDir: rootDir,
	}
}

func (e *EditorCommandExecutor) Run(command string, ctx glsp.Context, executor utils.Commander) error {
	err := e.tryRun(command, ctx, executor)
	if err != nil {
		ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeError,
			Message: "Command failed: " + command + ", reason: " + err.Error(),
		})
		return nil
	}

	return nil
}

func (e *EditorCommandExecutor) tryRun(command string, ctx glsp.Context, executor utils.Commander) error {
	e.mutex.Lock()
	v, found := e.commands[command]
	if !found {
		return errors.New("not found")
	}
	if v.running {
		return errors.New("already running")
	}
	v.running = true
	go v.fn(command, e, ctx, executor)
	e.mutex.Unlock()

	return nil
}

func (e *EditorCommandExecutor) resetRunning(command string) {
	e.mutex.Lock()
	if v, found := e.commands[command]; found {
		v.running = false
	}
	e.mutex.Unlock()
}
