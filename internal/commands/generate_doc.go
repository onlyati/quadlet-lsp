package commands

import (
	"encoding/json"
	"os"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/parser"
)

func generateDoc(
	command string,
	e *EditorCommandExecutor,
	messenger utils.Messenger,
	executor utils.Commander,
) {
	e.mutex.Lock()
	rootDir := e.rootDir
	e.mutex.Unlock()

	qd, err := parser.ParseQuadletDir(rootDir)
	if err != nil {
		messenger.SendMessage(
			utils.MessengerError,
			"failed to generate doc: "+err.Error(),
		)
	}

	file, err := json.Marshal(qd)
	if err != nil {
		messenger.SendMessage(
			utils.MessengerError,
			"failed to generate doc: "+err.Error(),
		)
	}
	os.WriteFile("./data.json", file, 0o644)

	messenger.SendMessage(
		utils.MessengerInfo,
		"Document is generated",
	)
}
