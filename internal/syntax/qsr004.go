package syntax

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Check if image name is fully qualified.
func qsr004(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	findings := utils.FindItems(
		s.documentText,
		"Container",
		"Image",
	)

	for _, findind := range findings {
		if strings.Count(findind.Value, "/") != 2 {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: findind.LineNumber, Character: 0},
					End:   protocol.Position{Line: findind.LineNumber, Character: findind.Length},
				},
				Severity: &s.warnDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr004"),
				Message:  "Image name is not fully qualified",
			})
		}
	}

	return diags
}
