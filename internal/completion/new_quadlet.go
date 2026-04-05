package completion

import (
	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listNewQuadletTemplates() []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	for k, category := range data.CategoryProperty {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:            k,
			Detail:           category.Details,
			InsertText:       category.InsertText,
			InsertTextFormat: &insertFormat,
			Kind:             &itemKind,
		})
	}

	return completionItems
}
