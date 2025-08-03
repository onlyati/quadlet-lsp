package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify name of Container, Pod, Network, Volume
func qsr020(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	var findings []utils.QuadletLine
	allowedFiles := []string{"container", "pod", "network", "volume"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			c+"Name",
		)
	}

	for _, finding := range findings {
		match := namingConvention.MatchString(finding.Value)
		if !match {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr020"),
				Message:  fmt.Sprintf("Invalid name of unit: %s", finding.Value),
			})
		}
	}

	return diags
}
