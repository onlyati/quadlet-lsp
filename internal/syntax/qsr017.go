package syntax

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Pod file does not exist
func qsr017(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		workspaceRoot := ""
		if s.config != nil {
			workspaceRoot = s.config.WorkspaceRoot
		}

		podFindings := utils.FindItems(
			utils.FindItemProperty{
				URI:           s.uri,
				RootDirectory: workspaceRoot,
				Text:          s.documentText,
				Section:       c,
				Property:      "Pod",
			},
		)

		for _, podFinding := range podFindings {
			podName := podFinding.Value
			if !strings.HasSuffix(podName, ".pod") {
				continue
			}

			// Extract directory from the URI
			uriPath := strings.TrimPrefix(s.uri, "file://")
			fileDir := path.Dir(uriPath)

			// Check if pod file exists adjacent to the container file
			adjacentPath := path.Join(fileDir, podName)
			_, err := os.Stat(adjacentPath)

			if errors.Is(err, os.ErrNotExist) {
				diags = append(diags, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: podFinding.LineNumber, Character: uint32(len(podFinding.Property) + 1)},
						End:   protocol.Position{Line: podFinding.LineNumber, Character: uint32(len(podFinding.Property) + 1 + len(podName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr017"),
					Message:  fmt.Sprintf("Pod file does not exists: %s", podName),
				})
			}
		}
	}

	return diags
}
