package syntax

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Network file does not exist
func qsr014(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "build", "kube"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Network",
		)
	}

	for _, finding := range findings {
		netName := finding.Value
		if strings.HasSuffix(netName, ".network") {
			_, err := os.Stat("./" + netName)

			if errors.Is(err, os.ErrNotExist) {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: finding.LineNumber, Character: uint32(len(finding.Property) + 1)},
						End:   protocol.Position{Line: finding.LineNumber, Character: uint32(len(finding.Property) + 1 + len(netName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr014"),
					Message:  fmt.Sprintf("Network file does not exists: %s", netName),
				})
				continue
			}

			if err != nil {
				log.Printf("failed to stat file: %s", err.Error())
				continue
			}
		}
	}

	return diags
}
