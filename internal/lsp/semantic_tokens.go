package lsp

import (
	"github.com/onlyati/quadlet-lsp/internal/semantic"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func SemanticTokens(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	uri := string(params.TextDocument.URI)
	text := docs.Read(uri)

	return semantic.CalculateSemanticTokens(text)
}
