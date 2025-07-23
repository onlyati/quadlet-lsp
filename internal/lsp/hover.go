package lsp

import (
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	uri := string(params.TextDocument.URI)
	text := documents.read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line

	section := findSection(lines, editorLine)
	if section == "" {
		return nil, nil
	}

	property := strings.Split(lines[editorLine], "=")[0]

	for _, item := range propertiesMap[section] {
		if property == item.label {
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**" + item.label + "**\n\n" + strings.Join(item.hover, "\n"),
				},
			}, nil
		}
	}

	return nil, nil
}
