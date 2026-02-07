package lsp

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func ExecuteCommands(context *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
	messenger := utils.ContextMessenger{
		Context: context,
	}
	err := commander.Run(params, &messenger, utils.CommandExecutor{})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
