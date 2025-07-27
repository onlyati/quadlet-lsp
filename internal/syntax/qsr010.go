package syntax

import (
	"strconv"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Validate that the PublishPort line is correct
func qsr010(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "kube"}
	var findings []utils.QuadletLine

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		findings = utils.FindItems(
			s.documentText,
			c,
			"PublishPort",
		)
	}

	for _, finding := range findings {
		tmp := strings.Split(finding.Value, ":")

		if len(tmp) == 2 {
			// This is something like PublishPort=420:69
			//                                   ^ this is the offset
			if tmp[0] != "" {
				tmpDiags := qsr010ValidatePortNumber(
					tmp[0],
					finding,
					len(finding.Property)+1,
				)
				diags = append(diags, tmpDiags...)
			}

			tmpDiags := qsr010ValidatePortNumber(
				tmp[1],
				finding,
				len(finding.Property)+1+len(tmp[0])+1,
			)
			diags = append(diags, tmpDiags...)

			continue
		}

		if len(tmp) == 3 {
			// This is something like PublishPort=127.0.0.1:420:69
			//                                             ^ this is the offset
			if tmp[1] != "" {
				tmpDiags := qsr010ValidatePortNumber(
					tmp[1],
					finding,
					len(finding.Property)+1+len(tmp[0])+1,
				)
				diags = append(diags, tmpDiags...)
			}

			tmpDiags := qsr010ValidatePortNumber(
				tmp[2],
				finding,
				len(finding.Property)+1+len(tmp[1])+1+len(tmp[0])+1,
			)
			diags = append(diags, tmpDiags...)

			continue
		}

		// If this is reached then line is probably incorrect format
		diags = append(diags, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: finding.LineNumber, Character: 0},
				End:   protocol.Position{Line: finding.LineNumber, Character: finding.Length},
			},
			Severity: &errDiag,
			Message:  "Incorrect format of PublishPort: invalid format",
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr010"),
		})
	}

	return diags
}

func qsr010ValidatePortNumber(text string, finding utils.QuadletLine, offset int) []protocol.Diagnostic {
	var diags []protocol.Diagnostic
	number, err := strconv.Atoi(text)
	if err != nil {
		diags = append(diags, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: finding.LineNumber, Character: uint32(offset)},
				End:   protocol.Position{Line: finding.LineNumber, Character: uint32(len(text) + offset)},
			},
			Severity: &errDiag,
			Message:  "Incorrect format of PublishPort: not a number",
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr010"),
		})
	}
	if number < 0 || number > 65535 {
		diags = append(diags, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: finding.LineNumber, Character: uint32(offset)},
				End:   protocol.Position{Line: finding.LineNumber, Character: uint32(len(text) + offset)},
			},
			Severity: &errDiag,
			Message:  "Incorrect format of PublishPort: port must be between [0;65535]",
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr010"),
		})
	}

	return diags
}
