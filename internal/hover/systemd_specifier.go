package hover

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleSystemSpecifier(info HoverInformation) *protocol.Hover {
	standOnPercentage := info.Line[info.CharacterPosition] == '%'
	if !standOnPercentage {
		return nil
	}

	// Before looking for the specifier it needs to be checked if
	// this is an escaped percentage sign '%%'. What is happening here
	// is to count '%' sings before the cursor, then decide it is escaped.
	if info.CharacterPosition > 0 {
		percentageCount := 0
		for i := info.CharacterPosition; ; i-- {
			if info.Line[i] == '%' {
				percentageCount += 1
			}

			if i == 0 || info.Line[i] != '%' {
				break
			}
		}
		if percentageCount%2 == 0 {
			return nil
		}
	}

	specifier := info.Line[info.CharacterPosition : info.CharacterPosition+2]
	specifierData, found := data.SystemdSpecifierSet[specifier]

	if !found {
		// The specifier does not exists, do nothing syntax checker is different module
		return nil
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "**" + specifierData.ShortDescription + "**\n\n" + strings.Join(specifierData.LongDescription, "\n"),
		},
	}
}
