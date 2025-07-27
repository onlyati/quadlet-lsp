package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify value of AutoUpdate
func qsr005(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "kube"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"AutoUpdate",
		)
	}

	if len(findings) > 0 {
		for _, f := range findings {
			if f.Value != "registry" && f.Value != "local" {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: f.LineNumber, Character: 0},
						End:   protocol.Position{Line: f.LineNumber, Character: f.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr005"),
					Message:  fmt.Sprintf("Invalid value of AutoUpdate: %s", f.Value),
				})
			}
		}
	}

	return diags
}
