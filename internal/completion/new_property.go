package completion

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listNewProperties(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	// If this is a continuation line, then it is not a new property
	if s.line > 0 {
		t := strings.TrimSpace(s.text[s.line-1])
		if strings.HasSuffix(t, "\\") {
			return completionItems
		}
	}

	// Normal processing
	s.config.Mu.RLock()
	podVer := s.config.Podman
	s.config.Mu.RUnlock()

	for _, p := range data.PropertiesMap[s.section] {
		checkVersion := podVer.GreaterOrEqual(p.MinVersion)
		if checkVersion {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: p.Label,
				Documentation: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**" + p.Label + "**\n\n" + strings.Join(p.Hover, "\n"),
				},
				Kind: &completionKind,
			})
		}
	}

	return completionItems
}
