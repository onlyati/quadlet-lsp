package completion

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func isItPropertyCompletion(line string) bool {
	return strings.Contains(line, "=")
}

func listPropertyCompletions(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	specialCompletions := map[string]func(s Completion) []protocol.CompletionItem{
		"Image":       propertyListImages,
		"Secret":      propertyListSecrets,
		"Volume":      propertyListVolumes,
		"Pod":         propertyListPods,
		"Network":     propertyListNetworks,
		"PublishPort": propertyListPorts,
		"UserNS":      propertyListUserIDs,
	}

	tmp := strings.Split(s.text[s.line], "=")
	propName := tmp[0]

	// Special handling on the line
	if v, ok := specialCompletions[propName]; ok {
		comps := v(s)
		if len(comps) > 0 {
			return comps
		}
	}

	// Generic suggestion based on data module
	s.config.Mu.RLock()
	podmanVer := s.config.Podman
	s.config.Mu.RUnlock()

	for _, p := range data.PropertiesMap[s.section] {
		labelCheck := propName == p.Label
		versionCheck := podmanVer.GreaterOrEqual(p.MinVersion)
		if labelCheck && versionCheck {
			for _, parm := range p.Parameters {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: parm,
				})
			}
		}
	}

	return completionItems
}

func propertyListPods(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	return completionItems
}

func propertyListNetworks(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	return completionItems
}

func propertyListPorts(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	return completionItems
}

func propertyListUserIDs(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	return completionItems
}
