package lsp

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// It is a very basic implementation. It is checking which section is,
// like `[Volume]`, `[Container]`, etc. then looking for the property that is
// in the current line. Then gather the document based on that and send the
// markdown response back.
func textHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	uri := string(params.TextDocument.URI)
	text := documents.Read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line

	section := findSection(lines, editorLine)
	if section == "" {
		return nil, nil
	}

	property := strings.Split(lines[editorLine], "=")[0]

	for _, item := range data.PropertiesMap[section] {
		if property == item.Label {
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**" + item.Label + "**\n\n" + strings.Join(item.Hover, "\n"),
				},
			}, nil
		}
	}

	return nil, nil
}
