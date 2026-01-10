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
			{Label: "rw", Kind: &valueKind},
			{Label: "ro", Kind: &valueKind},
			{Label: "z", Kind: &valueKind},
			{Label: "Z", Kind: &valueKind},
			{Label: "O", Kind: &valueKind},
			{Label: "copy", Kind: &valueKind},
			{Label: "nocopy", Kind: &valueKind},
			{Label: "dev", Kind: &valueKind},
			{Label: "nodev", Kind: &valueKind},
			{Label: "exec", Kind: &valueKind},
			{Label: "noexec", Kind: &valueKind},
			{Label: "suid", Kind: &valueKind},
			{Label: "nosuid", Kind: &valueKind},
			{Label: "bind", Kind: &valueKind},
			{Label: "rbind", Kind: &valueKind},
			{Label: "slave", Kind: &valueKind},
			{Label: "rslave", Kind: &valueKind},
			{Label: "shared", Kind: &valueKind},
			{Label: "rshared", Kind: &valueKind},
			{Label: "private", Kind: &valueKind},
			{Label: "rprivate", Kind: &valueKind},
			{Label: "unbindable", Kind: &valueKind},
			{Label: "runbindable", Kind: &valueKind},
		}
	}

	// Here we are after the '=' but before any ','
	// Suggest volumes from file and from the system
	volumes, err := utils.ListQuadletFiles("volume", s.config.WorkspaceRoot, *s.config.Project.DirLevel)
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
				Kind:  &valueKind,
			})
		}
	}

	return completionItems
}
