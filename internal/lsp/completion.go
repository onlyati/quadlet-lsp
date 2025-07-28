package lsp

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/completion"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// This function handles the completion event that is received.
func textCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	// executor := utils.CommandExecutor{}
	uri := string(params.TextDocument.URI)
	text := documents.Read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	charPos := params.Position.Character

	s := completion.NewCompletion(
		lines,
		uri,
		editorLine,
		charPos,
	)

	comps := s.RunCompletion(config)
	return comps, nil
}
