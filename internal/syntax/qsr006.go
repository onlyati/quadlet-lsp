package syntax

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr006AnotherQuadlet = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*\.(image|build)$`)

// Verify if the specified `.image` or `.build` file exists
// in the current working directory
func qsr006(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "volume"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Image"}: {},
			},
			qsr006Action,
		)
	}

	return diags
}

func qsr006Action(q utils.QuadletLine, _ utils.PodmanVersion) *protocol.Diagnostic {
	if !qsr006AnotherQuadlet.MatchString(q.Value) {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("failed to detect cwd: %s", err.Error())
		return nil
	}
	filePath := path.Join(cwd, q.Value)
	_, err = os.Stat(filePath)

	if !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return &protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: q.LineNumber, Character: 0},
			End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
		},
		Severity: &errDiag,
		Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr006"),
		Message:  fmt.Sprintf("Image file does not exists: %s", q.Value),
	}
}
