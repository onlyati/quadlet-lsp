package lsp

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/hover"
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
	cursorPos := params.Position.Character

	section := findSection(lines, editorLine)
	if section == "" {
		return nil, nil
	}

	config.Mu.RLock()
	rootDir := config.WorkspaceRoot
	level := config.Project.DirLevel
	config.Mu.RUnlock()

	return hover.HoverFunction(hover.HoverInformation{
		Line:              lines[editorLine],
		CharacterPosition: cursorPos,
		Section:           section,
		URI:               uri,
		LineNumber:        editorLine,
		RootDir:           rootDir,
		Level:             *level,
	}), nil
}
