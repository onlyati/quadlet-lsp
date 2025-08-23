package syntax

import (
	"fmt"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify name of Container, Pod, Network, Volume
func qsr020(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "network", "volume"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		propertyName := strings.TrimPrefix(c, "[")
		propertyName = strings.TrimSuffix(propertyName, "]")
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: propertyName + "Name"}: {},
			},
			qsr020Action,
		)
	}

	return diags
}

func qsr020Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	match := namingConvention.MatchString(q.Value)
	if match {
		return nil
	}
	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &errDiag,
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr020"),
			Message:  fmt.Sprintf("Invalid name of unit: %s", q.Value),
		},
	}
}
