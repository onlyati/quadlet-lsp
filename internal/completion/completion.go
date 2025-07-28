package completion

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	insertFormat = protocol.InsertTextFormatSnippet
	itemKind     = protocol.CompletionItemKindSnippet
)

type Completion struct {
	text      []string
	line      protocol.UInteger
	char      protocol.UInteger
	uri       string
	commander utils.Commander
	config    *utils.QuadletConfig
	section   string
}

func NewCompletion(
	document []string,
	uri string,
	currentLine, charPos protocol.UInteger,
) Completion {
	return Completion{
		text:      document,
		uri:       uri,
		commander: utils.CommandExecutor{},
		line:      currentLine,
		char:      charPos,
		section:   utils.FindSection(document, currentLine),
	}
}

func (s Completion) RunCompletion(config *utils.QuadletConfig) []protocol.CompletionItem {
	s.config = config
	var completionItems []protocol.CompletionItem

	// Section suggestions, things that are between '[]'
	if isItSectionLine(s.text[s.line]) {
		return listSections(s)
	}

	// If 'new.Something' is typed it provides suggestions for templates
	if isItNewMacro(s.text[s.line]) {
		return listNewMacros(s)
	}

	// There is a '=' in the line, so check for property's value
	if isItPropertyCompletion(s.text[s.line]) {
		return listPropertyCompletions(s)
	}

	return completionItems
}
