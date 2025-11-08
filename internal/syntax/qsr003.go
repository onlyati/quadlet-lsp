package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Checking for invalid properties
func qsr003(s SyntaxChecker) []protocol.Diagnostic {
	diags := []protocol.Diagnostic{}

	s.config.Mu.RLock()
	podmanVersion := s.config.Podman
	s.config.Mu.RUnlock()

	allowedFiles := []string{"image", "container", "volume", "network", "kube", "pod", "build", "artifact"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			podmanVersion,
			map[utils.ScanProperty]struct{}{
				{Section: "*", Property: "*"}: {},
			},
			qsr003Action,
		)
	}

	return diags
}

func qsr003Action(q utils.QuadletLine, p utils.PodmanVersion) []protocol.Diagnostic {
	section := q.Section[1 : len(q.Section)-1]
	if section == "Service" {
		// The [Service] is not implemented
		return nil
	}

	properties, foundSection := data.PropertiesMap[section]

	if !foundSection {
		// In this case we are in the line of the section header
		// Check that the section header exists
		return []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &errDiag,
				Message:  fmt.Sprintf("Invalid property is found: %s.%s", section, q.Property),
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr003"),
			},
		}
	}

	if q.Property == "" && q.Value == "" {
		// This only happen if it is a header line
		return nil
	}

	// The section exists the property should be checked now
	for _, prop := range properties {
		if prop.Label == q.Property && p.GreaterOrEqual(prop.MinVersion) {
			return nil
		}
	}

	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &errDiag,
			Message:  fmt.Sprintf("Invalid property is found: %s.%s", section, q.Property),
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr003"),
		},
	}
}
