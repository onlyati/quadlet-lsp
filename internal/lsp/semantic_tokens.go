package lsp

import (
	"github.com/onlyati/quadlet-lsp/internal/semantic"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func SemanticTokens(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	uri := string(params.TextDocument.URI)
	docs.ParseMutex.Lock() // Because the parse must be finished first, it must wait
	quadlet := docs.ReadLexerTokens(uri)
	docs.ParseMutex.Unlock()

	return semantic.CalculateSemanticTokens(quadlet)
}
