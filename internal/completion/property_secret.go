package completion

import (
	"log"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func propertyListSecrets(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	currentLine := s.text[s.line][:s.char]
	currentValue := strings.Split(currentLine, "=")[1]

	// Cursor is somewhere like `Secret=my-secret,`
	if strings.Contains(currentValue, ",") {
		return []protocol.CompletionItem{
			{
				Label: "type=mount",
			},
			{
				Label: "type=env",
			},
			{
				Label: "target=",
			},
		}
	}

	// Read the existing secrets and return with them
	output, err := s.commander.Run(
		"podman",
		"secret", "ls", "--format", "{{ .Name }}",
	)
	if err != nil {
		log.Println(err.Error())
	} else {
		for _, secret := range output {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: secret,
			})
		}
	}

	return completionItems
}
