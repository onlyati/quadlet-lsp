package syntax

import (
	"fmt"
	"strings"

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
		)
	}

	return diags
}

func qsr022Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	offset := uint32(len(q.Property)) + 1 // +1 is the '=' sign
	diags := []protocol.Diagnostic{}

	for path := range strings.SplitSeq(q.Value, ":") {
		if strings.Contains(path, "/%") {
			for i, c := range path {
				if c != '/' || i > len(path)-2 {
					continue
				}

				if path[i+1] != '%' || path[i+2] == ' ' {
					continue
				}

				specifier := path[i+1 : i+3]
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
		}
		offset += uint32(len(path)) + 1 // +1 is the ':' sign
	}

	return diags
}
