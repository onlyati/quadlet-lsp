// Package completion contains everything that is related for any completion
// from logical view.
package completion

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	insertFormat   = protocol.InsertTextFormatSnippet
	itemKind       = protocol.CompletionItemKindSnippet
	completionKind = protocol.CompletionItemKindKeyword
	valueKind      = protocol.CompletionItemKindValue
)

type Completion struct {
	text      []string
	line      protocol.UInteger
	char      protocol.UInteger
	uri       string
	commander utils.Commander
	config    *utils.QuadletConfig
	quadlet   *parser.QuadletNode
	tokenInfo parser.FindTokenOutput
}

func NewCompletion(
	document []string,
	uri string,
	currentLine, charPos protocol.UInteger,
	quadlet *parser.QuadletNode,
	tokenInfo parser.FindTokenOutput,
) Completion {
	return Completion{
		text:      document,
		uri:       uri,
		commander: utils.CommandExecutor{},
		line:      currentLine,
		char:      charPos,
		// section:   utils.FindSection(document, currentLine),
		quadlet:   quadlet,
		tokenInfo: tokenInfo,
	}
}

func (s Completion) RunCompletion(config *utils.QuadletConfig) []protocol.CompletionItem {
	s.config = config

	// File has no sections provide new templates
	if len(s.quadlet.Sections) == 0 {
		return listNewQuadletTemplates()
	}

	// Ignore comment
	if _, ok := s.tokenInfo.CurrentNode.(*parser.CommentNode); ok {
		return nil
	}

	// User start to type '[something'
	if _, ok := s.tokenInfo.CurrentNode.(*parser.SectionNode); ok {
		return listSections()
	}

	// Next conditions analyze the parent nodes, so if none return nothing
	if len(s.tokenInfo.ParentNodes) == 0 {
		return nil
	}

	// First parent is a section, user probably type a property
	if len(s.tokenInfo.ParentNodes) == 1 {
		if _, ok := s.tokenInfo.ParentNodes[0].(*parser.SectionNode); ok {
			return listNewProperties(s)
		}
	}

	// If user type '%' suggest systemd specifiers
	if s.line < uint32(len(s.text)) {
		if isItSystemSpecifier(s.text[s.line], s.char) {
			return listSystemdSpecifier(s)
		}
	}

	// If having two parents, then it must be a value position
	if len(s.tokenInfo.ParentNodes) == 2 {
		return listPropertyCompletions(s)
	}

	return nil

	//
	// // If user type '%' suggest systemd specifiers
	// if isItSystemSpecifier(s.text[s.line], s.char) {
	// 	return listSystemdSpecifier(s)
	// }
	//
	// // There is a '=' in the line, so check for property's value
	// if isItPropertyCompletion(s) {
	// 	return listPropertyCompletions(s)
	// }
	//
	// // If this point is reached, then user probably type something
	// // at the beginning of a file, let's suggest some property.
	// // Or typing before the '=' sign in a line
	// return listNewProperties(s)
}
