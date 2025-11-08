package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Checking for section headers.
func qsr001(s SyntaxChecker) []protocol.Diagnostic {
	allowedFiles := []string{"image", "container", "volume", "network", "kube", "pod", "build", "artifact"}
	if c := canFileBeApplied(s.uri, allowedFiles); c == "" {
		return []protocol.Diagnostic{}
	}

	units := []string{
		"[Image]",
		"[Container]",
		"[Volume]",
		"[Network]",
		"[Kube]",
		"[Pod]",
		"[Build]",
		"[Artifact]",
	}
	lines := strings.SplitSeq(s.documentText, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if slices.Contains(units, line) {
			return nil
		}
	}

	diag := protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: 0, Character: 0},
		},
		Severity: &errDiag,
		Message:  fmt.Sprintf("Missing any of these sections: %v", units),
		Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr001"),
	}

	return []protocol.Diagnostic{diag}
}
