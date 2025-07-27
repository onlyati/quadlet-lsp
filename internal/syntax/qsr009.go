package syntax

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify Annotation property
func qsr009(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "build", "network", "volume", "pod"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Label",
		)
	}

	for _, finding := range findings {
		index := strings.Index(finding.Value, "=")

		// Check if '=' is missing
		cond1 := index == -1

		// Cannot be space before or after the '=' sign
		cond2 := false
		if !cond1 {
			cond2 = finding.Value[index-1] == ' '

			if len(finding.Value)-1 > index {
				cond2 = cond2 || finding.Value[index+1] == ' '
			}
		}

		cond3 := index == len(finding.Value)-1

		if cond1 || cond2 || cond3 {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Message:  "Invalid format of Label specification",
				Severity: &s.errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr009"),
			})
		}
	}

	return diags
}
