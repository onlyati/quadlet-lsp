package lsp

import (
	"github.com/onlyati/quadlet-lsp/internal/hover"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// It is a very basic implementation. It is checking which section is,
// like `[Volume]`, `[Container]`, etc. then looking for the property that is
// in the current line. Then gather the document based on that and send the
// markdown response back.
func textHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	uri := string(params.TextDocument.URI)
	quadlet := docs.ReadQuadlet(uri)
	tokenInfo := quadlet.FindToken(
		parser.NodePosition{
			LineNumber: params.Position.Line,
			Position:   params.Position.Character,
		},
	)

	if len(tokenInfo.ParentNodes) == 0 {
		return nil, nil
	}

	config.Mu.RLock()
	rootDir := config.WorkspaceRoot
	level := config.Project.DirLevel
	config.Mu.RUnlock()

	return hover.HoverFunction(hover.HoverInformation{
		RootDir:   rootDir,
		Level:     *level,
		TokenInfo: tokenInfo,
	}), nil
}
