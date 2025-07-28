package completion

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listNewProperties(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	s.config.Mu.RLock()
	podVer := s.config.Podman
	s.config.Mu.RUnlock()

	for _, p := range data.PropertiesMap[s.section] {
		checkVersion := podVer.GreaterOrEqual(p.MinVersion)
		if checkVersion {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: p.Label + "=",
				Documentation: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**" + p.Label + "**\n\n" + strings.Join(p.Hover, "\n"),
				},
			})
		}
	}

	return completionItems
}
