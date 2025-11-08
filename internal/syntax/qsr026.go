package syntax

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify if artifcat is present in Artifact Quadlets
func qsr026(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"artifact"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		// Start the check only for *.artifact files, ignore dropins
		if strings.HasSuffix(s.uri, ".conf") {
			return diags
		}

		// Now scan everything (including dropins) and check if Artifact specified
		lines := utils.FindItems(utils.FindItemProperty{
			URI:           s.uri,
			RootDirectory: s.config.WorkspaceRoot,
			Text:          s.documentText,
			Section:       "[Artifact]",
			Property:      "Artifact",
		})

		if len(lines) == 0 {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 1},
				},
				Severity: &errDiag,
				Message:  "Artifact Quadlet file does not have Artifact property",
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr026"),
			},
			)
		}
	}

	return diags
}
