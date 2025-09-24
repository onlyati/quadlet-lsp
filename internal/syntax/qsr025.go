package syntax

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify if image is present in container Quadlets
func qsr025(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		// Start the check only for *.container files, ignore dropins
		if strings.HasSuffix(s.uri, ".conf") {
			return diags
		}

		// Now scan everything (including dropins) and check if Image specified
		lines := utils.FindItems(utils.FindItemProperty{
			URI:           s.uri,
			RootDirectory: s.config.WorkspaceRoot,
			Text:          s.documentText,
			Section:       "[Container]",
			Property:      "Image",
		})

		if len(lines) == 0 {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 1},
				},
				Severity: &errDiag,
				Message:  "Container Quadlet file does not have Image property",
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr025"),
			},
			)
		}
	}

	return diags
}
