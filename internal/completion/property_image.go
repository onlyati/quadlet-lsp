package completion

import (
	"log"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func propertyListImages(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	images, err := utils.ListQuadletFiles("*.image")
	if err != nil {
		log.Println(err.Error())
	}
	completionItems = append(completionItems, images...)

	builds, err := utils.ListQuadletFiles("*.build")
	if err != nil {
		log.Println(err.Error())
	}
	completionItems = append(completionItems, builds...)

	output, err := s.commander.Run(
		"podman",
		"images", "--format", "{{ .Repository }}:{{ .Tag }}",
	)
	if err != nil {
		log.Println(err.Error())
	} else {
		for _, image := range output {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: image,
			})
		}
	}

	return completionItems
}
