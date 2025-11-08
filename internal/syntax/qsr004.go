package syntax

import (
	"regexp"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr004FullyQualifiedImage = regexp.MustCompile(
	`^(?:[a-z0-9]+(?:[a-z0-9._-]+)*\.(?:[a-z0-9]+)|localhost)(?::[0-9]+)?(?:\/[a-z0-9]+(?:[a-zA-Z0-9-._]*))+`,
)

// Check if image name is fully qualified.
func qsr004(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "image", "volume", "artifact"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{}, // Does not matter just placeholder here
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Image"}:    {},
				{Section: c, Property: "Artifact"}: {},
			},
			qsr004Action,
		)
	}

	return diags
}

func qsr004Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	isItImage := strings.HasSuffix(q.Value, ".image")
	isItBuild := strings.HasSuffix(q.Value, ".build")
	if isItBuild || isItImage {
		return nil
	}
	if qsr004FullyQualifiedImage.MatchString(q.Value) {
		return nil
	}
	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &warnDiag,
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr004"),
			Message:  "Image name is not fully qualified",
		},
	}
}
