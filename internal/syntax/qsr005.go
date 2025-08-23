package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify value of AutoUpdate
func qsr005(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "kube"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{}, // placeholder
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "AutoUpdate"}: {},
			},
			qsr005Action,
		)
	}

	return diags
}

func qsr005Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	if q.Value == "registry" || q.Value == "local" {
		return nil
	}
	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &errDiag,
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr005"),
			Message:  fmt.Sprintf("Invalid value of AutoUpdate: %s", q.Value),
		},
	}
}
