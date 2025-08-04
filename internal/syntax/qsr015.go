package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr015ValidFlags = []string{
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

// Invalid format of Volume specification
func qsr015(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "build"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Volume"}: {},
			},
			qsr015Action,
		)
	}

	return diags
}

func qsr015Action(q utils.QuadletLine, _ utils.PodmanVersion) *protocol.Diagnostic {
	tmp := strings.Split(q.Value, ":")

	if len(tmp) >= 2 {
		if !strings.HasPrefix(tmp[1], "/") {
			return &protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr015"),
				Message:  "Invalid format of Volume specification: container directory is not absolute",
			}
		}
	}

	if len(tmp) == 3 {
		// Verify flags
		for f := range strings.SplitSeq(tmp[2], ",") {
			if !slices.Contains(qsr015ValidFlags, f) {
				return &protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: 0},
						End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr015"),
					Message:  fmt.Sprintf("Invalid format of Volume specification: '%s' flag is unknown", f),
				}
			}
		}
	}

	return nil
}
