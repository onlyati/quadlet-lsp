package completion

import (
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func isItSectionLine(line string) bool {
	return strings.HasPrefix(line, "[")
}

func listSections(_ Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	for k := range data.PropertiesMap {
		log.Println(k)
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: k,
		})
	}

	return completionItems
}
