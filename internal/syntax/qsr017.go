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

// Pod file does not exist
func qsr017(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Pod",
		)
	}

	for _, finding := range findings {
		podName := finding.Value
		if strings.HasSuffix(podName, ".pod") {
			_, err := os.Stat("./" + podName)

			if errors.Is(err, os.ErrNotExist) {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: finding.LineNumber, Character: uint32(len(finding.Property) + 1)},
						End:   protocol.Position{Line: finding.LineNumber, Character: uint32(len(finding.Property) + 1 + len(podName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr017"),
					Message:  fmt.Sprintf("Pod file does not exists: %s", podName),
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
