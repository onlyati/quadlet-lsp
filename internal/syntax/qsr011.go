package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	qsr011PortsRaw    []string
	qsr011FailedCheck []string
	qsr011Ports       []string
)

// The exposed port is not present in the image
func qsr011(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod"}
	var findigs []utils.QuadletLine
	qsr011PortsRaw = utils.FindImageExposedPorts(s.commander, s.uri)

	for _, s := range qsr011PortsRaw {
		if strings.HasPrefix(s, "failed-check-") {
			s, _ := strings.CutPrefix(s, "failed-check-")
			qsr011FailedCheck = append(qsr011FailedCheck, s)
		} else {
			if !slices.Contains(qsr011Ports, s) {
				qsr011Ports = append(qsr011Ports, s)
			}
		}
	}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "PublishPort"}: {},
			},
			qsr011Action,
		)
	}

	if len(findigs) == 0 {
		return diags
	}

	return diags
}

func qsr011Action(q utils.QuadletLine, _ utils.PodmanVersion) *protocol.Diagnostic {
	tmp := strings.Split(q.Value, ":")
	usedPort := tmp[len(tmp)-1]

	if !slices.Contains(qsr011Ports, usedPort) {
		if len(qsr011FailedCheck) == 0 {
			return &protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr011"),
				Message:  fmt.Sprintf("Port is not exposed in the image, exposed ports: %v", qsr011Ports),
			}
		} else {
			return &protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &infoDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr011"),
				Message:  fmt.Sprintf("Not able to verify exposed ports, because image not pulled: %v", qsr011FailedCheck),
			}
		}
	}

	return nil
}
