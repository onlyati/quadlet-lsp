package lsp

import (
	"errors"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/format"
	"github.com/onlyati/quadlet-lsp/internal/syntax"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Format(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	uri := string(params.TextDocument.URI)
	text := docs.Read(uri)
	textLines := strings.Split(text, "\n")

	// Only make formatting if no syntax error in the file
	checker := syntax.NewSyntaxChecker(text, uri)
	diags := checker.RunAll(config)
	if len(diags) > 0 {
		return nil, nil
	}

	p := docs.ReadParser(uri)
	if len(p.Errors) > 0 {
		return nil, errors.New("does not format until error present in file")
	}

	newText := format.FormatDocument(p.Quadlet)

	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: uint32(len(textLines)), Character: uint32(len(textLines[len(textLines)-1]))},
			},
			NewText: newText,
		},
	}, nil
}
