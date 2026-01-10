package completion

import (
	"log"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func propertyListPods(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	pods, err := utils.ListQuadletFiles("pod", s.config.WorkspaceRoot, *s.config.Project.DirLevel)
	if err != nil {
		log.Println(err.Error())
	} else {
		completionItems = append(completionItems, pods...)
	}

	return completionItems
}
