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

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Network"}: {},
			},
			qsr014Action,
		)
	}

	return diags
}

func qsr014Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	netName := q.Value
	if strings.HasSuffix(netName, ".network") {
		_, err := os.Stat("./" + netName)

		if errors.Is(err, os.ErrNotExist) {
			return []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1)},
						End:   protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1 + len(netName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr014"),
					Message:  fmt.Sprintf("Network file does not exists: %s", netName),
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
