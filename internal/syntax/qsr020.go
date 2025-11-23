package syntax

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr020NamingConvention = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

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
			nil,
		)
	}

	return diags
}

func qsr020Action(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
	match := qsr020NamingConvention.MatchString(q.Value)
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
