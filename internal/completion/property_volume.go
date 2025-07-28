package completion

import (
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func propertyListVolumes(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	currentLine := s.text[s.line][:s.char]
	currentValue := strings.Split(currentLine, "=")[1]

	// Example: Volume=foo.volume:/app/data:rw
	//                           ^         ^ After the 2nd specific flags can be
	//                           Delimiter between outside and inside location
	numberOfColons := strings.Count(currentValue, ":")

	if numberOfColons == 1 {
		// Do not give anything, typeing location inside the container
		return []protocol.CompletionItem{}
	}

	if numberOfColons == 2 {
		// Suggest some flag
		return []protocol.CompletionItem{
			{Label: "rw"},
			{Label: "ro"},
			{Label: "z"},
			{Label: "Z"},
		}
	}

	// Here we are after the '=' but before any ','
	// Suggest volumes from file and from the system
	volumes, err := utils.ListQuadletFiles("*.volume")
	if err != nil {
		log.Println(err.Error())
	} else {
		completionItems = append(completionItems, volumes...)
	}

	output, err := s.commander.Run(
		"podman",
		"volume", "ls", "--format", "{{ .Name }}",
	)
	if err != nil {
		log.Println(err.Error())
	} else {
		for _, volume := range output {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: volume,
			})
		}
	}

	return completionItems
}
