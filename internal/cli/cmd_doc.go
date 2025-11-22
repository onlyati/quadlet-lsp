package cli

import (
	"errors"
	"os"

	"github.com/onlyati/quadlet-lsp/internal/commands"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// CliMessenger command output normally send back to the editor. But when it is called
// from command line, there is no editor. So this structure implement the Messanger
// interface and catch the outputs from the command.
type CliMessenger struct {
	messages []struct {
		level utils.MessengerLevel
		text  string
	}
}

func (c *CliMessenger) SendMessage(level utils.MessengerLevel, text string) {
	c.messages = append(c.messages, struct {
		level utils.MessengerLevel
		text  string
	}{level, text})
}

// runDocCLI Run document generation from command line interface.
func runDocCLI(args []string, commander utils.Commander) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	parm := ""
	if len(args) >= 1 {
		parm = args[0]
	}
	e := commands.NewEditorCommandExecutor(cwd)
	messenger := CliMessenger{}
	commands.GenerateDoc(
		&protocol.ExecuteCommandParams{
			Command:   "generateDoc",
			Arguments: []any{parm},
		},
		&e,
		&messenger,
		commander,
	)

	output := []string{}
	for _, m := range messenger.messages {
		if m.level == utils.MessengerError || m.level == utils.MessengerWarning {
			err = errors.New("found warning or error messages")
		}
		output = append(output, m.text)
	}

	return output, err
}
