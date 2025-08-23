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

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Pod"}: {},
			},
			qsr017Action,
		)
	}

	return diags
}

func qsr017Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	podName := q.Value
	if strings.HasSuffix(podName, ".pod") {
		_, err := os.Stat("./" + podName)

		if errors.Is(err, os.ErrNotExist) {
			return []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1)},
						End:   protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1 + len(podName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr017"),
					Message:  fmt.Sprintf("Pod file does not exists: %s", podName),
				},
			}
		}

		if err != nil {
			log.Printf("failed to stat file: %s", err.Error())
			return nil
		}
	}

	return nil
}
