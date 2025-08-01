package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Invalid format of Volume specification
func qsr015(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "build"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Volume",
		)
	}

	validFlags := []string{
		"rw",
		"ro",
		"z",
		"Z",
		"O",
		"U",
		"copy", "nocopy",
		"dev", "nodev",
		"exec", "noexec",
		"suid", "nosuid",
		"bind", "rbind",
		"shared", "shared",
		"slave", "rslave",
		"private", "rprivate",
		"unbindable", "runbindable",
	}

	for _, finding := range findings {
		tmp := strings.Split(finding.Value, ":")

		if len(tmp) >= 2 {
			if !strings.HasPrefix(tmp[1], "/") {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: finding.LineNumber, Character: 0},
						End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr015"),
					Message:  "Invalid format of Volume specification: container directory is not absolute",
				})
			}
		}

		if len(tmp) == 3 {
			// Verify flags
			for f := range strings.SplitSeq(tmp[2], ",") {
				if !slices.Contains(validFlags, f) {
					diags = append(diags, protocol.Diagnostic{
						Range: protocol.Range{
							Start: protocol.Position{Line: finding.LineNumber, Character: 0},
							End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
						},
						Severity: &errDiag,
						Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr015"),
						Message:  fmt.Sprintf("Invalid format of Volume specification: '%s' flag is unknown", f),
					})
				}
			}
		}
	}

	return diags
}
