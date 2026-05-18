package lsp

import (
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/completion"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// This function handles the completion event that is received.
func textCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	log.Println("received a completion request")
	// executor := utils.CommandExecutor{}
	uri := string(params.TextDocument.URI)
	text := docs.Read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	charPos := params.Position.Character

	quadlet := docs.ReadQuadlet(uri)
	tokenInfo := quadlet.FindToken(
		parser.NodePosition{
			LineNumber: params.Position.Line,
			Position:   params.Position.Character,
		},
	)

	s := completion.NewCompletion(
		lines,
		uri,
		editorLine,
		charPos,
		&quadlet,
		tokenInfo,
	)
	comps := s.RunCompletion(config)

	log.Printf("return %d completion", len(comps))
	return comps, nil
}
