package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var qsr015ValidFlags = []string{
	"rw",
	"ro",
	"z",
	"Z",
	"O",
	"U",
	"copy", "nocopy",
	"dev", "nodev",
	"exec", "noexec",
	"suid", "nosuid",
	"bind", "rbind",
	"shared", "shared",
	"slave", "rslave",
	"private", "rprivate",
	"unbindable", "runbindable",
}

// Invalid format of Volume specification
func qsr015(s SyntaxChecker) []protocol.Diagnostic {
	var diags []protocol.Diagnostic

	allowedFiles := []string{"container", "pod", "build"}

	if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
		diags = utils.ScanQadlet(
			s.documentText,
			utils.PodmanVersion{},
			map[utils.ScanProperty]struct{}{
				{Section: c, Property: "Volume"}: {},
			},
			qsr015Action,
		)
	}

	return diags
}

func qsr015Action(q utils.QuadletLine, _ utils.PodmanVersion) *protocol.Diagnostic {
	tmp := strings.Split(q.Value, ":")

	if len(tmp) >= 2 {
		containerPath := tmp[1]

		// If path start with '%' then user may want to use systemd specifier
		// like %h instead of home. It has to be checked this is a directory.
		startsWithValidSpecifier := false
		if strings.HasPrefix(containerPath, "%") && len(containerPath) > 1 {
			specifier := containerPath[0:2]
			data, found := data.SystemdSpecifierSet[specifier]
			if found && data.IsDirectory {
				startsWithValidSpecifier = true
			}
		}

		starsWithSlash := strings.HasPrefix(containerPath, "/")
		if !starsWithSlash && !startsWithValidSpecifier {
			return &protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: q.LineNumber, Character: 0},
					End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
				},
				Severity: &errDiag,
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr015"),
				Message:  "Invalid format of Volume specification: container directory is not absolute",
			}
		}
	}

	if len(tmp) == 3 {
		// Verify flags
		flags := tmp[2]
		for f := range strings.SplitSeq(flags, ",") {
			if !slices.Contains(qsr015ValidFlags, f) {
				return &protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: q.LineNumber, Character: 0},
						End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
					},
					Severity: &errDiag,
					Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr015"),
					Message:  fmt.Sprintf("Invalid format of Volume specification: '%s' flag is unknown", f),
				}
			}
		}
	}

	return nil
}
