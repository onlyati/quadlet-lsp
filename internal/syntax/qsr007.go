package syntax

import (
	"fmt"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Check syntax of Environment property
func qsr007(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "build", "build"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Environment"}: {},
			},
			qsr007Action,
		)
	}

	return diags
}

func qsr007Action(q utils.QuadletLine, _ utils.PodmanVersion) *protocol.Diagnostic {
	tokens, err := splitQuoted(q.Value)
	if err != nil {
		return &protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &errDiag,
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr007"),
			Message:  fmt.Sprintf("Invalid format: %s", err.Error()),
		}
	}

	invalids := []string{}
	for _, token := range tokens {
		if quotedKV.MatchString(token) || unquotedKV.MatchString(token) || aposthropeKV.MatchString(token) {
			continue
		}
		invalids = append(invalids, token)
	}

	if len(invalids) == 0 {
		return nil
	}

	return &protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: q.LineNumber, Character: 0},
			End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
		},
		Severity: &errDiag,
		Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr007"),
		Message:  fmt.Sprintf("Invalid format: bad delimiter usage at %s", strings.Join(invalids, ", ")),
	}
}
