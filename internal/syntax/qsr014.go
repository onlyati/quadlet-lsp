package syntax

import (
	"fmt"
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type qsr014ActionParms struct {
	rootDir  string
	dirLevel int
}

// Network file does not exist
func qsr014(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "build", "kube"}
	s.config.Mu.RLock()
	parm := qsr014ActionParms{
		rootDir:  s.config.WorkspaceRoot,
		dirLevel: *s.config.Project.DirLevel,
	}
	s.config.Mu.RUnlock()

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Network"}: {},
			},
			qsr014Action,
			parm,
		)
	}

	return diags
}

func qsr014Action(q utils.QuadletLine, _ utils.PodmanVersion, extraInfo any) []protocol.Diagnostic {
	actionParm := qsr014ActionParms{}
	switch v := extraInfo.(type) {
	case qsr014ActionParms:
		actionParm = v
	default:
		return nil
	}

	netName := q.Value
	if strings.HasSuffix(netName, ".network") {
		quadlets, err := utils.ListQuadletFiles("network", actionParm.rootDir, actionParm.dirLevel)
		exists := false

		for _, q := range quadlets {
			if netName == q.Label {
				exists = true
				break
			}
		}

		if !exists {
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
