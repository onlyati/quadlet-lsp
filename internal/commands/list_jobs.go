package commands

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listJobs(command *protocol.ExecuteCommandParams, e *EditorCommandExecutor, messenger utils.Messenger, executor utils.Commander) {
	defer e.resetRunning(command.Command)

	e.mutex.Lock()
	runningTasks := []string{}
	for k, v := range e.commands {
		if v.running {
			runningTasks = append(runningTasks, k)
		}
	}
	e.mutex.Unlock()

	messenger.SendMessage(
		utils.MessengerInfo,
		fmt.Sprintf("Running tasks: %+v", runningTasks),
	)
}
