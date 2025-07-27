package syntax

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Check if image name is fully qualified.
func qsr004(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "image", "volume"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Image",
		)
	}

	for _, finding := range findings {
		if strings.HasSuffix(finding.Value, ".image") || strings.HasSuffix(finding.Value, ".build") {
			continue
		}
		if strings.Count(finding.Value, "/") != 2 {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Severity: &warnDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr004"),
				Message:  "Image name is not fully qualified",
			})
		}
	}

	return diags
}
