package syntax

import (
	"fmt"
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type qsr017ActionParms struct {
	rootDir  string
	dirLevel int
}

// Pod file does not exist
func qsr017(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container"}
	s.config.Mu.RLock()
	parm := qsr017ActionParms{
		rootDir:  s.config.WorkspaceRoot,
		dirLevel: *s.config.Project.DirLevel,
	}
	s.config.Mu.RUnlock()

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Pod"}: {},
			},
			qsr017Action,
			parm,
		)
	}

	return diags
}

func qsr017Action(q utils.QuadletLine, _ utils.PodmanVersion, extraInfo any) []protocol.Diagnostic {
	actionParm := qsr017ActionParms{}
	switch v := extraInfo.(type) {
	case qsr017ActionParms:
		actionParm = v
	default:
		return nil
	}

	podName := q.Value
	if strings.HasSuffix(podName, ".pod") {
		quadlets, err := utils.ListQuadletFiles("pod", actionParm.rootDir, actionParm.dirLevel)
		exists := false

		for _, q := range quadlets {
			if podName == q.Label {
				exists = true
				break
			}
		}

		if !exists {
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
