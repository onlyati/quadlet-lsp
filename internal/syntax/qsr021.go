package syntax

import (
	"fmt"
	"regexp"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	qsr021ServiceNamingConvention = regexp.MustCompile(`^[a-zA-Z0-9][@a-zA-Z0-9_.-]*\.service$`)
	qsr021QuadletNamingConvention = regexp.MustCompile(
		`^[a-zA-Z0-9][@a-zA-Z0-9_.-]*\.(image|container|volume|network|kube|pod|build)$`,
	)
)

// Wrong depdency format is used
func qsr021(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic
	allowedFiles := []string{"image", "container", "volume", "network", "kube", "pod", "build"}

	s.config.Mu.RLock()
	podmanVer := s.config.Podman
	s.config.Mu.RUnlock()

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			podmanVer,
			map[utils.ScanProperty]struct{}{
				{Section: "[Unit]", Property: "Wants"}:     {},
				{Section: "[Unit]", Property: "Requires"}:  {},
				{Section: "[Unit]", Property: "Requisite"}: {},
				{Section: "[Unit]", Property: "BindsTo"}:   {},
				{Section: "[Unit]", Property: "PartOf"}:    {},
				{Section: "[Unit]", Property: "Upholds"}:   {},
				{Section: "[Unit]", Property: "Conflicts"}: {},
				{Section: "[Unit]", Property: "Before"}:    {},
				{Section: "[Unit]", Property: "After"}:     {},
			},
			qsr021Action,
		)
	}

	return diags
}

func qsr021Action(q utils.QuadletLine, p utils.PodmanVersion) []protocol.Diagnostic {
	if qsr021ServiceNamingConvention.MatchString(q.Value) {
		return nil
	}

	if p.GreaterOrEqual(utils.BuildPodmanVersion(5, 5, 0)) {
		if qsr021QuadletNamingConvention.MatchString(q.Value) {
			return nil
		}
	}
	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: q.LineNumber, Character: 0},
				End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
			},
			Severity: &errDiag,
			Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr021"),
			Message:  fmt.Sprintf("Invalid depdency is specified: %s", q.Value),
		},
	}
}
