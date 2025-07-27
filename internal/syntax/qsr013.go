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

// Volume file does not exists
func qsr013(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "build"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Volume",
		)
	}

	for _, finding := range findings {
		tmp := strings.Split(finding.Value, ":")
		if len(tmp) == 0 {
			continue
		}

		volName := tmp[0]
		if strings.HasSuffix(volName, ".volume") {
			_, err := os.Stat("./" + volName)

			if errors.Is(err, os.ErrNotExist) {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: finding.LineNumber, Character: uint32(len(finding.Property) + 1)},
						End:   protocol.Position{Line: finding.LineNumber, Character: uint32(len(finding.Property) + 1 + len(volName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr013"),
					Message:  fmt.Sprintf("Volume file does not exists: %s", volName),
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
