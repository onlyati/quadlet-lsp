package completion

import (
	"log"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func propertyListNetworks(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	volumes, err := utils.ListQuadletFiles("network", s.config.WorkspaceRoot, *s.config.Project.DirLevel)
	if err != nil {
		log.Println(err.Error())
	} else {
		completionItems = append(completionItems, volumes...)
	}

	output, err := s.commander.Run(
		"podman",
		"network", "ls", "--format", "{{ .Name }}",
	)
	if err != nil {
		log.Println(err.Error())
	} else {
		for _, network := range output {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: network,
				Kind:  &valueKind,
			})
		}
	}

	return completionItems
}
