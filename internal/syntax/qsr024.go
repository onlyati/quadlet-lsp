package syntax

import (
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Usage on User, Group and DynamicUser in rootless Podman is not recommended
func qsr024(s SyntaxChecker) []protocol.Diagnostic {
	diags := []protocol.Diagnostic{}

	allowedFiles := []string{"image", "container", "volume", "network", "kube", "pod", "build"}
	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: "[Service]", Property: "User"}:        {},
				{Section: "[Service]", Property: "Group"}:       {},
				{Section: "[Service]", Property: "DynamicUser"}: {},
			},
			qsr024Action,
			nil,
		)
	}

	return diags
}

func qsr024Action(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &warnDiag,
			Message:  fmt.Sprintf("Usage in rootless podman is not recommended: %s.%s", "Service", q.Property),
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr024"),
		},
	}
}
