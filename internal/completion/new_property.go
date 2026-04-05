package completion

import (
	"fmt"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func listNewProperties(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	parentToken, ok := s.tokenInfo.ParentNodes[0].(*parser.SectionNode)
	if !ok {
		return nil
	}

	// Normal processing
	s.config.Mu.RLock()
	podVer := s.config.Podman
	s.config.Mu.RUnlock()

	section := strings.TrimPrefix(*parentToken.Text, "[")
	section = strings.TrimSuffix(section, "]")
	for _, p := range data.PropertiesMap[section] {
		checkVersion := podVer.GreaterOrEqual(p.MinVersion)
		var textEdit protocol.TextEdit
		if checkVersion {
			if p.Macro != "" {
				textEdit = protocol.TextEdit{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      s.line,
							Character: 0,
						},
					},
					NewText: p.Macro,
				}
			} else {
				textEdit = protocol.TextEdit{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      s.line,
							Character: 0,
						},
					},
					NewText: fmt.Sprintf("%s=${1:value}\n$0", p.Label),
				}
			}
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: p.Label,
				Documentation: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**" + p.Label + "**\n\n" + strings.Join(p.Hover, "\n"),
				},
				Kind:             &itemKind,
				TextEdit:         textEdit,
				InsertTextFormat: &insertFormat,
			})
		}
	}

	return completionItems
}
