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

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Volume"}: {},
			},
			qsr013Action,
		)
	}

	return diags
}

func qsr013Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	tmp := strings.Split(q.Value, ":")
	if len(tmp) == 0 {
		return nil
	}

	volName := tmp[0]
	if strings.HasSuffix(volName, ".volume") {
		if strings.Contains(volName, "@") {
			volName = utils.ConvertTemplateNameToFile(volName)
		}
		_, err := os.Stat("./" + volName)

		if errors.Is(err, os.ErrNotExist) {
			return []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1)},
						End:   protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1 + len(volName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr013"),
					Message:  fmt.Sprintf("Volume file does not exists: %s", volName),
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
