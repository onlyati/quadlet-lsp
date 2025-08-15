package commands

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func listJobs(command string, e *EditorCommandExecutor, messenger utils.Messenger, executor utils.Commander) {
	defer e.resetRunning(command)

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
