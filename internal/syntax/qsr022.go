package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Check if path at volume has '/%x' characters and 'x' is not directory
func qsr022(s SyntaxChecker) []protocol.Diagnostic {
	diags := []protocol.Diagnostic{}

	allowedFiles := []string{"container", "pod", "build"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Volume"}: {},
			},
			qsr022Action,
			nil,
		)
	}

	return diags
}

func qsr022Action(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
	offset := uint32(len(q.Property)) + 1 // +1 is the '=' sign
	diags := []protocol.Diagnostic{}

	for i, c := range q.Value {
		if c != '/' || i > len(q.Value)-3 {
			continue
		}

		if q.Value[i+1] != '%' || q.Value[i+2] == ' ' {
			continue
		}

		specifier := q.Value[i+1 : i+3]
		data, found := data.SystemdSpecifierSet[specifier]

		if !found {
			// Non existing verifier is used, check in different syntax rule
			continue
		}

		if data.IsDirectory {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: offset + uint32(i)},
					End:   protocol.Position{Line: q.LineNumber, Character: offset + uint32(i) + 2},
				},
				Severity: &errDiag,
				Message:  fmt.Sprintf("Specifier, %s, already starts with '/' sign", specifier),
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr022"),
			})
		}
	}

	return diags
}
