package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify for any invalid systemd specifier
func qsr023(s SyntaxChecker) []protocol.Diagnostic {
	diags := []protocol.Diagnostic{}

	allowedFiles := []string{"image", "container", "volume", "network", "kube", "pod", "build"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: "*", Property: "*"}: {},
			},
			qsr023Action,
			nil,
		)
	}
	return diags
}

func qsr023Action(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
	offset := uint32(len(q.Property)) + 1 // +1 is the '=' sign
	diags := []protocol.Diagnostic{}

	for i, c := range q.Value {
		if c != '%' || i > len(q.Value)-2 {
			continue
		}

		if q.Value[i+1] == ' ' {
			continue
		}

		specifier := q.Value[i : i+2]
		_, found := data.SystemdSpecifierSet[specifier]

		if !found {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: offset + uint32(i)},
					End:   protocol.Position{Line: q.LineNumber, Character: offset + uint32(i) + 1},
				},
				Severity: &errDiag,
				Message:  fmt.Sprintf("Specifier '%s' is invalid", specifier),
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr023"),
			})
		}
	}

	return diags
}
