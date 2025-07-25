package syntax

import (
	"fmt"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify value of AutoUpdate
func qsr005(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	if !strings.HasSuffix(s.uri, ".container") && strings.HasSuffix(s.uri, ".kube") {
		return diags
	}

	tmp := strings.Split(s.uri, ".")
	section := utils.FirstCharacterToUpper(tmp[len(tmp)-1])

	findings := utils.FindItems(
		s.documentText,
		section,
		"AutoUpdate",
	)

	if len(findings) > 0 {
		for _, f := range findings {
			if f.Value != "registry" && f.Value != "local" {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: f.LineNumber, Character: 0},
						End:   protocol.Position{Line: f.LineNumber, Character: f.Length},
					},
					Severity: &s.errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr005"),
					Message:  fmt.Sprintf("Invalid value of AutoUpdate: %s", f.Value),
				})
			}
		}
	}

	return diags
}
