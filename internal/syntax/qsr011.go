package syntax

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	qsr011PortsRaw    []string
	qsr011FailedCheck []string
	qsr011Ports       []string
	qsr011Mutex       sync.Mutex
)

// The exposed port is not present in the image
func qsr011(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	qsr011Mutex.Lock()
	defer qsr011Mutex.Unlock()
	qsr011PortsRaw = nil
	qsr011FailedCheck = nil
	qsr011Ports = nil

	allowedFiles := []string{"container", "pod"}
	qsr011PortsRaw = utils.FindImageExposedPorts(
		s.commander,
		s.uri,
		s.config.WorkspaceRoot,
		s.uri,
	)

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
		items := utils.FindItems(
			utils.FindItemProperty{
				URI:           s.uri,
				RootDirectory: s.config.WorkspaceRoot,
				Text:          s.documentText,
				Section:       c,
				Property:      "PublishPort",
			},
		)

		for _, q := range items {
			tempDiags := qsr011Action(q, utils.PodmanVersion{})
			diags = append(diags, tempDiags...)
		}
	}

	return diags
}

func qsr011Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
	tmp := strings.Split(q.Value, ":")
	usedPort := tmp[len(tmp)-1]

	tmp = strings.Split(usedPort, "/")
	usedPort = tmp[0]

	if !slices.Contains(qsr011Ports, usedPort) {
		if len(qsr011FailedCheck) == 0 {
			return []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: 0},
						End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr011"),
					Message:  fmt.Sprintf("Port is not exposed in the image, exposed ports: %v", qsr011Ports),
				},
			}
		} else {
			return []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: 0},
						End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
					},
					Severity: &infoDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr011"),
					Message:  fmt.Sprintf("Not able to verify exposed ports, because image not pulled: %v", qsr011FailedCheck),
				},
			}
		}
	}

	return nil
}
