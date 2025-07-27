package syntax

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify if the specified `.image` or `.build` file exists
// in the current working directory
func qsr006(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "volume"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Image",
		)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("failed to detect cwd: %s", err.Error())
		return diags
	}

	for _, finding := range findings {
		if !strings.HasSuffix(finding.Value, ".image") && !strings.HasSuffix(finding.Value, ".build") {
			continue
		}

		filePath := path.Join(cwd, finding.Value)
		_, err := os.Stat(filePath)

		if errors.Is(err, os.ErrNotExist) {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: finding.LineNumber, Character: 0},
					End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr006"),
				Message:  fmt.Sprintf("Image file does not exists: %s", finding.Value),
			})
		}
	}

	return diags
}
