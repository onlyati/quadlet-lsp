// Package hover
//
// This package contains hover actions and reactions for them.
package hover

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type HoverInformation struct {
	RootDir   string
	TokenInfo parser.FindTokenOutput
	Level     int
}

func HoverFunction(info HoverInformation) *protocol.Hover {
	switch info.TokenInfo.CurrentNode.(type) {
	case *parser.AssignNode:
		return handlePropertyHover(info)
	case *parser.ValueNode:
		return handleValueHover(info)
	}
	return nil
}

func handlePropertyHover(info HoverInformation) *protocol.Hover {
	if len(info.TokenInfo.ParentNodes) == 0 {
		return nil
	}
	section, ok := info.TokenInfo.ParentNodes[0].(*parser.SectionNode)
	if !ok {
		return nil
	}
	keyword, ok := info.TokenInfo.CurrentNode.(*parser.AssignNode)
	if !ok {
		return nil
	}

	sectionValue := strings.TrimPrefix(*section.Text, "[")
	sectionValue = strings.TrimSuffix(sectionValue, "]")
	for _, item := range data.PropertiesMap[sectionValue] {
		if *keyword.Name == item.Label {
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
		"UserNS":  handleValueUserNS,
		"Volume":  handleValueVolume,
		"Secret":  handleValueSecret,
		"Pod":     handleValuePod,
		"Network": handleValueNetwork,
	}
	keyword, ok := info.TokenInfo.ParentNodes[0].(*parser.AssignNode)
	if !ok {
		return nil
	}

	fn, found := handlerMap[*keyword.Name]
	if found {
		return fn(info)
	}

	return nil
}
