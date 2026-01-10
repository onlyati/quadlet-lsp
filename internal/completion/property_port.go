package completion

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func propertyListPorts(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	// Let's find out that we need to provide any complation at all
	colons := strings.Count(s.text[s.line], ":")
	tmp := strings.Split(s.text[s.line], ":")

	// We need complation in two cases:
	// ExposedPorts=127.0.0.1:420:69
	// ExposedPorts=420:69
	if colons == 0 {
		return completionItems
	}
	if colons == 1 {
		// Check if first part is an IP address
		if strings.Count(tmp[0], ".") > 0 {
			return completionItems
		}
	}

	// Now gather ports
	props := utils.FindImageExposedPortsProperty{
		C:        s.commander,
		URI:      s.uri,
		RootDir:  s.config.WorkspaceRoot,
		Name:     s.uri,
		DirLevel: *s.config.Project.DirLevel,
	}
	ports := utils.FindImageExposedPorts(props)
	for _, port := range ports {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: port,
			Kind:  &valueKind,
		})
	}

	return completionItems
}
