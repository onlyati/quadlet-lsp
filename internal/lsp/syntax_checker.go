package lsp

import (
	"github.com/onlyati/quadlet-lsp/internal/syntax"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func SyntaxCheckOnSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	uri := string(params.TextDocument.URI)
	text := documents.read(uri)

	checker := syntax.NewSyntaxChecker(text, uri)

	// Run all syntax checker rule
	diags := checker.RunAll(config)
	if len(diags) > 0 {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         protocol.DocumentUri(uri),
			Diagnostics: diags,
		})
	} else {
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         protocol.DocumentUri(uri),
			Diagnostics: []protocol.Diagnostic{},
		})
	}

	return nil
}
