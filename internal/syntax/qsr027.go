package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Verify value of ImageVolume
func qsr027(s SyntaxChecker) []protocol.Diagnostic {
	diags := []protocol.Diagnostic{}

	s.config.Mu.RLock()
	podmanVer := s.config.Podman
	s.config.Mu.RUnlock()

	allowedFiles := []string{"container"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			podmanVer,
			map[utils.ScanProperty]struct{}{
				{Section: "[Container]", Property: "ImageVolume"}: {},
			},
			qsr027Action,
			nil,
		)
	}

	return diags
}

func qsr027Action(q utils.QuadletLine, p utils.PodmanVersion, _ any) []protocol.Diagnostic {
	if !p.GreaterOrEqual(utils.BuildPodmanVersion(6, 1, 0)) {
		return nil
	}
	if q.Value != "bind" && q.Value != "tmpfs" && q.Value != "ignore" {
		return []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &errDiag,
				Message:  fmt.Sprintf("Allowed values are bind, tmpfs and ignore, it is: %s", q.Value),
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr027"),
			},
		}
	}

	return nil
}
