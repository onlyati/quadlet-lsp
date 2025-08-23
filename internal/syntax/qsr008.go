package syntax

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr008KeyValueCheck = regexp.MustCompile(`^(['"]?)([A-Za-z0-9][A-Za-z0-9.-_/]*)=(.*)(['"]?)$`)

// Verify Annotation property
func qsr008(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "build", "build"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Annotation"}: {},
			},
			qsr008Action,
		)
	}

	return diags
}

func qsr008Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	tokens, err := splitQuoted(q.Value)
	if err != nil {
		return []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr008"),
				Message:  fmt.Sprintf("Invalid format: %s", err.Error()),
			},
		}
	}

	invalids := []string{}
	for _, token := range tokens {
		if qsr008KeyValueCheck.MatchString(token) {
			continue
		}
		invalids = append(invalids, token)
	}

	if len(invalids) == 0 {
		return nil
	}

	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &errDiag,
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr008"),
			Message:  fmt.Sprintf("Invalid format: bad delimiter usage at %s", strings.Join(invalids, ", ")),
		},
	}
}
