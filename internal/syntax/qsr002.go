package syntax

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Checking for unfinihsed lines
func qsr002(s SyntaxChecker) []protocol.Diagnostic {
	allowedFiles := []string{"image", "container", "volume", "network", "kube", "pod", "build", "artifact"}
	if c := canFileBeApplied(s.uri, allowedFiles); c == "" {
		return []protocol.Diagnostic{}
	}

	var diags []protocol.Diagnostic

	lineNum := uint32(0)
	lines := strings.SplitSeq(s.documentText, "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "=") && strings.Count(line, "=") == 1 {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: lineNum, Character: 0},
					End:   protocol.Position{Line: lineNum, Character: uint32(len(line))},
				},
				Severity: &errDiag,
				Message:  "Line is unfinished",
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr002"),
			})
		}
		lineNum++
	}

	return diags
}
