package completion

import (
	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listSections() []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	for k := range data.PropertiesMap {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: k,
			Kind:  &completionKind,
		})
	}

	return completionItems
}
