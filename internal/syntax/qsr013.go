package syntax

import (
	"fmt"
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type qsr013ActionParms struct {
	rootDir  string
	dirLevel int
}

// Volume file does not exists
func qsr013(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "build"}
	s.config.Mu.RLock()
	parm := qsr013ActionParms{
		rootDir:  s.config.WorkspaceRoot,
		dirLevel: *s.config.Project.DirLevel,
	}
	s.config.Mu.RUnlock()

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Volume"}: {},
			},
			qsr013Action,
			parm,
		)
	}

	return diags
}

func qsr013Action(q utils.QuadletLine, _ utils.PodmanVersion, extraInfo any) []protocol.Diagnostic {
	tmp := strings.Split(q.Value, ":")
	if len(tmp) == 0 {
		return nil
	}

	actionParm := qsr013ActionParms{}
	switch v := extraInfo.(type) {
	case qsr013ActionParms:
		actionParm = v
	default:
		return nil
	}

	volName := tmp[0]
	if strings.HasSuffix(volName, ".volume") {
		if strings.Contains(volName, "@") {
			volName = utils.ConvertTemplateNameToFile(volName)
		}
		quadlets, err := utils.ListQuadletFiles("volume", actionParm.rootDir, actionParm.dirLevel)
		exists := false

		for _, q := range quadlets {
			if volName == q.Label {
				exists = true
				break
			}
		}

		if !exists {
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
