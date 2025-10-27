package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"text/template"

	"github.com/onlyati/quadlet-lsp/internal/embeds"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func generateDoc(
	command *protocol.ExecuteCommandParams,
	e *EditorCommandExecutor,
	messenger utils.Messenger,
	executor utils.Commander,
) {
	defer e.resetRunning(command.Command)
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

	if len(command.Arguments) > 0 {
		outType := command.Arguments[0]
		switch v := outType.(type) {
		case string:
			switch v {
			case "json":
				generateDocJSON(qd, rootDir, messenger)
			case "md":
				generateDocMd(qd, rootDir, messenger)
			default:
				messenger.SendMessage(
					utils.MessengerError,
					"failed to generate doc: expected parameters: json, md, html",
				)
			}
		default:
			messenger.SendMessage(
				utils.MessengerError,
				"failed to generate doc: expect string parameter",
			)
		}
	} else {
		messenger.SendMessage(
			utils.MessengerError,
			"failed to generate doc: no output type is specified",
		)
	}
}

func generateDocJSON(qd parser.QuadletDirectory, rootDir string, messenger utils.Messenger) {
	file, err := json.Marshal(qd)
	if err != nil {
		messenger.SendMessage(
			utils.MessengerError,
			"failed to generate doc: "+err.Error(),
		)
	}
	os.MkdirAll(path.Join(rootDir, "doc"), os.FileMode(0o755))
	os.WriteFile(path.Join(rootDir, "doc", "doc.json"), file, 0o644)

	messenger.SendMessage(
		utils.MessengerInfo,
		"Document is generated",
	)
}

func generateDocMd(qd parser.QuadletDirectory, rootDir string, messenger utils.Messenger) {
	t, err := template.ParseFS(embeds.TemplateFs, "*.tpl")
	if err != nil {
		messenger.SendMessage(
			utils.MessengerError,
			"Internal error: "+err.Error(),
		)
		return
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "md_main.tpl", qd)
	if err != nil {
		messenger.SendMessage(
			utils.MessengerError,
			"Internal error: "+err.Error(),
		)
		return
	}

	os.MkdirAll(path.Join(rootDir, "doc"), os.FileMode(0o755))
	os.WriteFile(path.Join(rootDir, "doc", "doc.md"), buf.Bytes(), 0o644)

	messenger.SendMessage(
		utils.MessengerInfo,
		"Document is generated",
	)
}
