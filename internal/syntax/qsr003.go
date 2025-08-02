package syntax

import (
	"fmt"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Checking for invalid properties
func qsr003(s SyntaxChecker) []protocol.Diagnostic {
	allowedFiles := []string{"image", "container", "volume", "network", "kube", "pod"}
	if c := canFileBeApplied(s.uri, allowedFiles); c == "" {
		return []protocol.Diagnostic{}
	}

	var diags []protocol.Diagnostic

	lineNum := uint32(0)
	section := ""
	props := []data.PropertyMapItem{}
	lines := strings.SplitSeq(s.documentText, "\n")

	s.config.Mu.RLock()
	podmanVersion := s.config.Podman
	s.config.Mu.RUnlock()

	for line := range lines {
		lineNum++
		line = strings.TrimSpace(line)

		// Skip emptry or comment lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Read the current section
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			tSection := line[1 : len(line)-1]

			if tProps, ok := data.PropertiesMap[tSection]; ok {
				section = tSection
				props = tProps
			} else {
				section = ""
				props = nil
			}
			continue
		}

		// If we are in a section then check for property
		if section != "" && len(props) > 0 && strings.Contains(line, "=") {
			tmp := strings.Split(line, "=")
			if len(tmp) == 0 {
				continue
			}

			found := false
			for _, prop := range props {
				if tmp[0] == prop.Label && podmanVersion.GreaterOrEqual(prop.MinVersion) {
					found = true
					break
				}
			}

			if found {
				continue
			}

			// If this point reached, then property not found
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{Line: lineNum - 1, Character: 0},
					End:   protocol.Position{Line: lineNum - 1, Character: uint32(len(line) - 1)},
				},
				Severity: &errDiag,
				Message:  fmt.Sprintf("Invalid property is found: %s.%s", section, tmp[0]),
				Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr003"),
			})
		}

	}

	return diags
}
