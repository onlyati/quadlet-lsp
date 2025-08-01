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
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"Secret",
		)
	}

	for _, finding := range findings {
		tmp := strings.Split(finding.Value, ",")

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
				diags = append(diags, qsr012MakeDiag(
					fmt.Sprintf(
						"Invalid format of secret specification: '%s' has no value",
						tmpPart[0],
					),
					finding,
				))
			} else {
				// Check the option name
				secretParms[tmpPart[0]] = tmpPart[1]
			}
		}

		for k, v := range secretParms {
			if k == "type" {
				if v != "mount" && v != "env" {
					diags = append(diags, qsr012MakeDiag(
						"Invalid format of secret specification: 'type' can be either 'mount' or 'env'",
						finding,
					))
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
					diags = append(diags, qsr012MakeDiag(
						fmt.Sprintf(
							"Invalid format of secret specification: '%s' only allowed if type=mount",
							k,
						),
						finding,
					))
				}
				continue
			}

			diags = append(diags, qsr012MakeDiag(
				fmt.Sprintf(
					"Invalid format of secret specification: '%s' is invalid option",
					k,
				),
				finding,
			))
		}
	}

	return diags
}

func qsr012MakeDiag(text string, finding utils.QuadletLine) protocol.Diagnostic {
	return protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{Line: finding.LineNumber, Character: 0},
			End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
		},
		Severity: &errDiag,
		Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr012"),
		Message:  text,
	}
}
