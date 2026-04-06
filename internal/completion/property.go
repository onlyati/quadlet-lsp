package completion

import (
	"log"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listPropertyCompletions(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	sectionNode, ok := s.tokenInfo.ParentNodes[1].(*parser.SectionNode)
	if !ok {
		return nil
	}
	parentNode, ok := s.tokenInfo.ParentNodes[0].(*parser.AssignNode)
	if !ok {
		return nil
	}

	specialCompletions := map[string]func(s Completion) []protocol.CompletionItem{
		"Image":       propertyListImages,
		"Secret":      propertyListSecrets,
		"Volume":      propertyListVolumes,
		"Pod":         propertyListPods,
		"Network":     propertyListNetworks,
		"PublishPort": propertyListPorts,
		"UserNS":      propertyListUserIDs,
	}

	propName := *parentNode.Name

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

	section := strings.TrimPrefix(*sectionNode.Text, "[")
	section = strings.TrimSuffix(section, "]")
	log.Println("section is: " + section)
	for _, p := range data.PropertiesMap[section] {
		labelCheck := propName == p.Label
		versionCheck := podmanVer.GreaterOrEqual(p.MinVersion)
		if labelCheck && versionCheck {
			for _, parm := range p.Parameters {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: parm,
					Kind:  &completionKind,
				})
			}
		}
	}

	return completionItems
}
