package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// The exposed port is not present in the image
func qsr011(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod"}
	var findigs []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findigs = utils.FindItems(
			s.documentText,
			c,
			"PublishPort",
		)
	}

	if len(findigs) == 0 {
		return diags
	}

	ports := utils.FindImageExposedPorts(s.commander, s.uri)

	for _, finding := range findigs {
		tmp := strings.Split(finding.Value, ":")
		usedPort := tmp[len(tmp)-1]

		if !slices.Contains(ports, usedPort) {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr011"),
				Message:  fmt.Sprintf("Port is not exposed in the image, exposed ports: %v", ports),
			})
		}
	}

	return diags
}
