package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Container cannot have PublishPort if belongs to a pod
func qsr018(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		publishFindings := utils.FindItems(
			s.documentText,
			c,
			"PublishPort",
		)
		podFindings := utils.FindItems(
			s.documentText,
			c,
			"Pod",
		)

		if len(publishFindings) > 0 && len(podFindings) > 0 {
			for _, p := range publishFindings {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: p.LineNumber, Character: 0},
						End:   protocol.Position{Line: p.LineNumber, Character: p.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr018"),
					Message:  fmt.Sprintf("Container cannot have PublishPort because belongs to a pod: %s", podFindings[0].Value),
				})
			}
		}
	}

	return diags
}
