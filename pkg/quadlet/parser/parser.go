// Package parser is a module to hold logic for Quadlet file parsing by tokens.
package parser

import (
	"os"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/lexer"
)

type ParserError struct {
	Text  string
	Token *lexer.Token
}

type Parser struct {
	Quadlet       *QuadletNode   // The parsed output
	Errors        []ParserError  // If there is any error during reading record it
	Path          string         // Location of the file
	tokens        []lexer.Token  // Result of the lexer
	tokenIndex    int            // Which token we are in iteration
	commentBuffer []*CommentNode // Comments before statements are belongs to statement
}

func NewParser(path string) Parser {
	content, err := os.ReadFile(path)
	if err != nil {
		content = []byte{}
	}
	l := lexer.NewLexer(string(content))
	l.Run()
	return Parser{
		Path: path,
		Quadlet: &QuadletNode{
			Documents: []*CommentNode{},
			Sections:  []*SectionNode{},
		},
		tokenIndex:    -1,
		tokens:        l.Tokens,
		commentBuffer: []*CommentNode{},
		Errors:        []ParserError{},
	}
}

// Run method iterate over the tokens and create an AST.
func (p *Parser) Run() {
	token := p.consumeToken()
	for token != nil {
		if token.Type != lexer.TokenTypeComment && len(p.Quadlet.Documents) == 0 {
			prevToken := p.prevToken()
			if prevToken != nil {
				if prevToken.Type == lexer.TokenTypeComment {
					p.Quadlet.Documents = p.commentBuffer
					p.commentBuffer = []*CommentNode{}
				}
			}
		}

		switch token.Type {
		case lexer.TokenTypeComment:
			p.parseComment(token)
		case lexer.TokenTypeSection:
			p.parseSection(token)
		case lexer.TokenTypeKeyword:
			p.parseAssignment(token)
		case lexer.TokenTypeValue, lexer.TokenTypeAssign, lexer.TokenTypeContSign:
			// The values are processed after Keyword so this should never hit
			// If this is hit, it is an unexpected token
			syntaxError := ParserError{
				Token: token,
				Text:  "unexpected token",
			}
			p.Errors = append(p.Errors, syntaxError)
		}

		token = p.consumeToken()
	}
}

func (p *Parser) parseComment(token *lexer.Token) {
	if len(p.commentBuffer) > 0 && len(p.Quadlet.Documents) == 0 {
		// If there was an empty line, then it is a Quadlet document
		// and not belongs to section
		if token.StartPos.LineNumber-p.commentBuffer[len(p.commentBuffer)-1].EndPos.LineNumber > 1 {
			p.Quadlet.Documents = p.commentBuffer
			p.commentBuffer = []*CommentNode{}
		}
	}

	// Just buffer the comment and as soon we hit a different token
	// we will assign it there as document.
	comment := CommentNode{
		StartPos: NodePosition(token.StartPos),
		EndPos:   NodePosition(token.EndPos),
		Text:     &token.Text,
	}
	p.commentBuffer = append(p.commentBuffer, &comment)
}

func (p *Parser) parseSection(token *lexer.Token) {
	section := SectionNode{
		StartPos:    NodePosition(token.StartPos),
		EndPos:      NodePosition(token.EndPos),
		Text:        &token.Text,
		Documents:   p.commentBuffer,
		Assignments: []*AssignNode{},
	}
	p.Quadlet.Sections = append(p.Quadlet.Sections, &section)
	p.commentBuffer = []*CommentNode{}
}

func (p *Parser) parseAssignment(token *lexer.Token) {
	assignment := AssignNode{
		StartPos:  NodePosition(token.StartPos),
		Documents: p.commentBuffer,
		Name:      &token.Text,
	}
	p.commentBuffer = []*CommentNode{}

	// If not in section it is an error
	if p.lastSection() == nil {
		p.Errors = append(p.Errors, ParserError{
			Text:  "keyword without section is invalid",
			Token: token,
		})
		return
	}

	// Next token should be an assign
	if p.peekToken() == nil {
		p.Errors = append(p.Errors, ParserError{
			Text:  "expects an '=' sign after keyword, it got end of file",
			Token: token,
		})
		return
	}
	if p.peekToken().Type != lexer.TokenTypeAssign {
		p.Errors = append(p.Errors, ParserError{
			Text:  "expects an '=' sign after keyword, it got " + p.peekToken().Text,
			Token: token,
		})
		return
	}
	token = p.consumeToken()

	// Read until we read value and continue sign
	value := ValueNode{}
	valueString := strings.Builder{}
	startNotSet := true
	if p.peekToken() == nil {
		p.Errors = append(p.Errors, ParserError{
			Text:  "unfinished line",
			Token: token,
		})
		return
	}
	for p.peekToken().Type == lexer.TokenTypeValue || p.peekToken().Type == lexer.TokenTypeContSign {
		token = p.consumeToken()

		if token.Type == lexer.TokenTypeValue {
			if startNotSet {
				value.StartPos = NodePosition(token.StartPos)
				startNotSet = false
			}
			valueString.WriteString(strings.TrimSpace(token.Text))
			valueString.WriteString(" ")
		}

		if p.peekToken() == nil {
			break
		}
	}

	value.EndPos = NodePosition(token.EndPos)
	value.Value = utils.AsPtr(strings.TrimSpace(valueString.String()))

	assignment.Value = &value
	assignment.EndPos = NodePosition(token.EndPos)
	p.lastSection().Assignments = append(p.lastSection().Assignments, &assignment)
}

func (p *Parser) consumeToken() *lexer.Token {
	p.tokenIndex++
	if p.tokenIndex == len(p.tokens)-1 {
		return nil
	}
	return &p.tokens[p.tokenIndex]
}

func (p *Parser) prevToken() *lexer.Token {
	if p.tokenIndex == 0 {
		return nil
	}
	return &p.tokens[p.tokenIndex-1]
}

func (p *Parser) peekToken() *lexer.Token {
	if p.tokenIndex+1 < len(p.tokens)-1 {
		return &p.tokens[p.tokenIndex+1]
	}
	return nil
}

func (p *Parser) lastSection() *SectionNode {
	if len(p.Quadlet.Sections) == 0 {
		return nil
	}
	return p.Quadlet.Sections[len(p.Quadlet.Sections)-1]
}
