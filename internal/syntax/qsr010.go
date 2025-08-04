package syntax

import (
	"regexp"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr010PortRegexp = regexp.MustCompile(
	`^([.|\d]*?\:?)(([1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))\:(([1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))$`,
)

// Validate that the PublishPort line is correct
func qsr010(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "kube"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText, utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "PublishPort"}: {},
			},
			qsr010Action,
		)
	}

	return diags
}

func qsr010Action(q utils.QuadletLine, _ utils.PodmanVersion) *protocol.Diagnostic {
	if qsr010PortRegexp.MatchString(q.Value) {
		return nil
	}

	return &protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: q.LineNumber, Character: 0},
			End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
		},
		Severity: &errDiag,
		Message:  "Incorrect format of PublishPort",
		Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr010"),
	}
}
