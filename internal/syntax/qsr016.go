package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Invalid value of UserNS specification
func qsr016(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "kube"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"UserNS",
		)
	}

	allowdUserNs := []string{
		"auto",
		"host",
		"keep-id",
		"nomap",
	}

	for _, finding := range findings {
		tmp := strings.Split(finding.Value, ":")

		if len(tmp) == 0 {
			continue
		}

		if len(tmp) > 1 {
			if tmp[0] != "keep-id" {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: finding.LineNumber, Character: 0},
						End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr016"),
					Message:  fmt.Sprintf("Invalid value of UserNS: '%s' has no parameters", tmp[0]),
				})
			} else {
				for p := range strings.SplitSeq(tmp[1], ",") {
					checkUid := strings.HasPrefix(p, "uid=")
					checkGid := strings.HasPrefix(p, "gid=")
					if !checkUid && !checkGid {
						diags = append(diags, protocol.Diagnostic{
							Range: protocol.Range{
								Start: protocol.Position{Line: finding.LineNumber, Character: 0},
								End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
							},
							Severity: &errDiag,
							Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr016"),
							Message:  fmt.Sprintf("Invalid value of UserNS: [uid gid] allowed but found %s", p),
						})
					}
				}
			}
		}

		if !slices.Contains(allowdUserNs, tmp[0]) {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr016"),
				Message:  fmt.Sprintf("Invalid value of UserNS: allowed values: '%v' and found %s", allowdUserNs, tmp[0]),
			})
		}
	}

	return diags
}
