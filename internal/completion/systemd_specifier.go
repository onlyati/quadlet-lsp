package completion

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func isItSystemSpecifier(line string, charPos protocol.UInteger) bool {
	if len(line) < 2 {
		return false
	}
	if line[charPos-1] != '%' {
		return false
	}

	// Count how many '%' sign before the cursor, if even then provide completion
	percentageCount := 0
	for i := charPos - 1; ; i-- {
		if i == 0 {
			break
		}

		if line[i] != '%' {
			break
		}

		percentageCount++
	}
	return percentageCount%2 == 1
}

func listSystemdSoecifier(s Completion) []protocol.CompletionItem {
	var completions []protocol.CompletionItem

	for k, v := range data.SystemdSpecifierSet {
		textEdit := protocol.TextEdit{
			Range: protocol.Range{
				Start: protocol.Position{Line: s.line, Character: s.char - 1},
				End:   protocol.Position{Line: s.line, Character: s.char + 1},
			},
			NewText: k,
		}
		completions = append(completions, protocol.CompletionItem{
			Label:            k,
			Kind:             &itemKind,
			TextEdit:         textEdit,
			InsertTextFormat: &insertFormat,
			Documentation: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**" + v.ShortDescription + "**\n\n" + strings.Join(v.LongDescription, "\n"),
			},
		})
	}

	return completions
}
