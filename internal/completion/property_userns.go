package completion

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// If somebody type `UserNS=keep-id:`, then check if image has any user
// defined, and provide its id for uid and gid as well
func propertyListUserIDs(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	currentValue := strings.Split(s.text[s.line], "=")[1]

	// Only interested if keep-id is used
	if !strings.HasPrefix(currentValue, "keep-id:") {
		return []protocol.CompletionItem{}
	}

	tmp := strings.Split(s.uri, ".")
	ext := tmp[len(tmp)-1]
	findings := utils.FindItems(
		utils.FindItemProperty{
			URI:           s.uri,
			RootDirectory: s.config.WorkspaceRoot,
			Text:          strings.Join(s.text, "\n"),
			Section:       "[" + utils.FirstCharacterToUpper(ext) + "]",
			Property:      "Image",
			DirLevel:      *s.config.Project.DirLevel,
		},
	)

	if len(findings) == 0 {
		return []protocol.CompletionItem{}
	}

	imageName := findings[0].Value
	output, err := s.commander.Run(
		"podman",
		"image", "inspect", imageName,
	)
	if err != nil {
		return nil
	}
	inspectJSON := strings.Join(output, "")
	log.Println(inspectJSON)
	var data []map[string]any
	err = json.Unmarshal([]byte(inspectJSON), &data)
	if err != nil {
		log.Printf("failed to unmarshal json: %s", err.Error())
	}

	if len(data) == 0 {
		log.Printf("image is not pulled: %s", imageName)
	}

	config, ok := data[0]["Config"].(map[string]any)
	if !ok {
		return nil
	}

	user, ok := config["User"].(string)
	if !ok {
		return nil
	}

	completionItems = append(completionItems, protocol.CompletionItem{
		Label: "uid=" + user,
		Kind:  &valueKind,
	})
	completionItems = append(completionItems, protocol.CompletionItem{
		Label: "gid=" + user,
		Kind:  &valueKind,
	})

	return completionItems
}
