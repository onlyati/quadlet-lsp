package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr016AllowdUserNs = []string{
	"auto",
	"host",
	"keep-id",
	"nomap",
}

// Invalid value of UserNS specification
func qsr016(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "kube"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText, utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "UserNS"}: {},
			},
			qsr016Action,
			nil,
		)
	}

	return diags
}

func qsr016Action(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
	tmp := strings.Split(q.Value, ":")

	if len(tmp) == 0 {
		return nil
	}

	if len(tmp) > 1 {
		if tmp[0] != "keep-id" {
			return []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: 0},
						End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr016"),
					Message:  fmt.Sprintf("Invalid value of UserNS: '%s' has no parameters", tmp[0]),
				},
			}
		} else {
			for p := range strings.SplitSeq(tmp[1], ",") {
				checkUID := strings.HasPrefix(p, "uid=")
				checkGID := strings.HasPrefix(p, "gid=")
				if !checkUID && !checkGID {
					return []protocol.Diagnostic{
						{
							Range: protocol.Range{
								Start: protocol.Position{Line: q.LineNumber, Character: 0},
								End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
							},
							Severity: &errDiag,
							Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr016"),
							Message:  fmt.Sprintf("Invalid value of UserNS: [uid gid] allowed but found %s", p),
						},
					}
				}
			}
		}
	}

	if !slices.Contains(qsr016AllowdUserNs, tmp[0]) {
		return []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr016"),
				Message:  fmt.Sprintf("Invalid value of UserNS: allowed values: '%v' and found %s", qsr016AllowdUserNs, tmp[0]),
			},
		}
	}

	return nil
}
