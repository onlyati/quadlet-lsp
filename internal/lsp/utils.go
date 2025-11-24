package lsp

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// This function looking for that the cursor currently in which section.
// Sections are like `[Container]`, `[Unit]`, and so on.
func findSection(lines []string, lineNumber protocol.UInteger) string {
	section := ""
	for i := lineNumber; ; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			break
		}

		if i == 0 {
			break
		}
	}
	return section
}
