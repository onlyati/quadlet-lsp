package hover

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type HoverInformation struct {
	Line              string
	Uri               string
	Section           string
	LineNumber        protocol.UInteger
	CharacterPosition protocol.UInteger
	property          string
	value             string
}

func HoverFunction(info HoverInformation) *protocol.Hover {
	splitLine := strings.SplitN(info.Line, "=", 2)
	if len(splitLine) < 2 {
		return nil
	}

	info.property = splitLine[0]
	info.value = splitLine[1]

	// Verify where the cursor is
	separatorPosition := strings.Index(info.Line, "=")
	if info.CharacterPosition < protocol.UInteger(separatorPosition) {
		return handlePropertyHover(info)
	} else {
		// Check if cursor on systemd specifier
		hoverValue := handleSystemSpecifier(info)
		if hoverValue != nil {
			return hoverValue
		}

		// Handle value specific hovers
		return handleValueHover(info)
	}
}

func handlePropertyHover(info HoverInformation) *protocol.Hover {
	for _, item := range data.PropertiesMap[info.Section] {
		if info.property == item.Label {
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**" + item.Label + "**\n\n" + strings.Join(item.Hover, "\n"),
				},
			}
		}
	}

	return nil
}

func handleValueHover(info HoverInformation) *protocol.Hover {
	handlerMap := map[string]func(HoverInformation) *protocol.Hover{
		"UserNS": handleValueUserNS,
		"Volume": handleValueVolume,
	}

	fn, found := handlerMap[info.property]
	if found {
		return fn(info)
	}

	return nil
}
