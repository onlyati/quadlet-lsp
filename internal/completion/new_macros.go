package completion

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func isItNewMacro(line string) bool {
	return strings.HasPrefix(line, "new.")
}

func listNewMacros(s Completion) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	propName := strings.TrimPrefix(s.text[s.line], "new.")

	s.config.Mu.RLock()
	podVer := s.config.Podman
	s.config.Mu.RUnlock()

	for _, p := range data.PropertiesMap[s.section] {
		versionCheck := podVer.GreaterOrEqual(p.MinVersion)
		prefixCheck := strings.HasPrefix(p.Label, propName)
		macroCheck := p.Macro != ""
		if versionCheck && prefixCheck && macroCheck {
			textEdit := protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      s.line,
						Character: 0,
					},
					End: protocol.Position{
						Line:      s.line,
						Character: uint32(len(s.text[s.line])),
					},
				},
				NewText: p.Macro,
			}

			completionItems = append(completionItems, protocol.CompletionItem{
				Label:            "new." + p.Label,
				Kind:             &itemKind,
				TextEdit:         textEdit,
				InsertTextFormat: &insertFormat,
			})
		}
	}

	return completionItems
}
