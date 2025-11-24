package syntax

import (
	"fmt"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Invalid format of Secret specification
func qsr012(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Secret"}: {},
			},
			qsr012Action,
			nil,
		)
	}

	return diags
}

func qsr012Action(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
	tmp := strings.Split(q.Value, ",")

	secretParms := map[string]string{}

	for i, part := range tmp {
		// Only interested in line like: Secret=my-secret,type=env,target=MYVAL
		// Skip of only secret name is specified
		if i == 0 {
			continue
		}

		tmpPart := strings.Split(part, "=")
		if len(tmpPart) == 1 {
			// Option value did not specified
			return qsr012MakeDiag(
				fmt.Sprintf(
					"'%s' has no value",
					tmpPart[0],
				),
				q,
			)
		} else {
			// Check the option name
			secretParms[tmpPart[0]] = tmpPart[1]
		}
	}

	for k, v := range secretParms {
		if k == "type" {
			if v != "mount" && v != "env" {
				return qsr012MakeDiag(
					"'type' can be either 'mount' or 'env'",
					q,
				)
			}
			continue
		}

		if k == "target" {
			// Nothing to check but it is a valid option
			continue
		}

		if k == "uid" || k == "gid" || k == "mode" {
			// They are invalid if the type is not mount
			// If type is omitted then default is mount
			// So they only invalid of type is "env"
			if secretParms["type"] == "env" {
				return qsr012MakeDiag(
					fmt.Sprintf(
						"'%s' only allowed if type=mount",
						k,
					),
					q,
				)
			}
			continue
		}

		return qsr012MakeDiag(
			fmt.Sprintf(
				"'%s' is invalid option",
				k,
			),
			q,
		)
	}

	return nil
}

func qsr012MakeDiag(text string, finding utils.QuadletLine) []protocol.Diagnostic {
	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: finding.LineNumber, Character: 0},
				End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
			},
			Severity: &errDiag,
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr012"),
			Message:  fmt.Sprintf("Invalid format of secret specification: %s", text),
		},
	}
}
