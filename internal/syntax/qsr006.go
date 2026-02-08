package syntax

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr006AnotherQuadlet = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*\.(image|build)$`)

type qsr006ActionParms struct {
	rootDir  string
	dirLevel int
}

// Verify if the specified `.image` or `.build` file exists
// in the current working directory
func qsr006(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "volume"}
	s.config.Mu.RLock()
	parm := qsr006ActionParms{
		rootDir:  s.config.WorkspaceRoot,
		dirLevel: *s.config.Project.DirLevel,
	}
	s.config.Mu.RUnlock()

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Image"}: {},
			},
			qsr006Action,
			parm,
		)
	}

	return diags
}

func qsr006Action(q utils.QuadletLine, _ utils.PodmanVersion, extraInfo any) []protocol.Diagnostic {
	actionParm := qsr006ActionParms{}
	switch v := extraInfo.(type) {
	case qsr006ActionParms:
		actionParm = v
	default:
		return nil
	}

	imageName := q.Value
	ext := q.Value[strings.LastIndexAny(q.Value, ".")+1:]
	if qsr006AnotherQuadlet.MatchString(q.Value) {
		quadlets, err := utils.ListQuadletFiles(ext, actionParm.rootDir, actionParm.dirLevel)
		exists := false

		for _, q := range quadlets {
			if imageName == q.Label {
				exists = true
				break
			}
		}

		if !exists {
			return []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1)},
						End:   protocol.Position{Line: q.LineNumber, Character: uint32(len(q.Property) + 1 + len(imageName))},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr006"),
					Message:  fmt.Sprintf("Image file does not exists: %s", imageName),
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
