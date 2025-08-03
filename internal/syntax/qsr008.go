package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify Annotation property
func qsr008(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "build"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Annotation",
		)
	}

	for _, finding := range findings {
		tokens, err := splitQuoted(finding.Value)
		if err != nil {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr008"),
				Message:  fmt.Sprintf("Invalid format: %s", err.Error()),
			})
			continue
		}

		for _, token := range tokens {
			if quotedKV.MatchString(token) || unquotedKV.MatchString(token) || aposthropeKV.MatchString(token) {
				continue
			}

			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr008"),
				Message:  fmt.Sprintf("Invalid format: bad delimiter usage at %s", token),
			})
		}
	}

	return diags
}
