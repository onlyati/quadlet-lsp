package commands

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listJobs(command string, e *EditorCommandExecutor, ctx glsp.Context, executor utils.Commander) {
	defer e.resetRunning(command)

	e.mutex.Lock()
	runningTasks := []string{}
	for k, v := range e.commands {
		if v.running {
			runningTasks = append(runningTasks, k)
		}
	}
	e.mutex.Unlock()

	ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: fmt.Sprintf("Running tasks: %+v", runningTasks),
	})
}
