// Package semantic contains functions, structs and constants that is related
// for semantic tokens.
package semantic

import (
	"slices"
	"unicode/utf8"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

var TokenLegends = []string{
	string(protocol.SemanticTokenTypeComment),   // Comment lines
	string(protocol.SemanticTokenTypeProperty),  // Like 'Image=', 'Exec=', 'Pod='
	string(protocol.SemanticTokenTypeNamespace), // Section like '[Container]', '[Unit]'
	string(protocol.SemanticTokenTypeString),    // Value belongs to keywords
	string(protocol.SemanticTokenTypeOperator),  // Operators like '=', ':', ','
	string(protocol.SemanticTokenTypeClass),     // Used within values to highlight things
	string(protocol.SemanticTokenTypeParameter), // Used within values to highlight things
}

func CalculateSemanticTokens(fileText string) (*protocol.SemanticTokens, error) {
	tokens := parseQuadlet(fileText)

	var data []uint32
	lastLine := uint32(0)
	lastChar := uint32(0)

	for _, token := range tokens {
		deltaLine := token.line - lastLine
		deltaStart := token.charPos

		if deltaLine == 0 {
			deltaStart = token.charPos - lastChar
		}

		data = append(data,
			deltaLine,
			deltaStart,
			token.length,
			uint32(slices.Index(TokenLegends, token.tokenType)),
			0, // No modifiers
		)

		lastLine = token.line
		lastChar = token.charPos
	}

	return &protocol.SemanticTokens{Data: data}, nil
}

type token struct {
	line      protocol.UInteger
	charPos   protocol.UInteger
	length    protocol.UInteger
	tokenType string
}

type lexer struct {
	input        string
	position     protocol.UInteger
	readPosition protocol.UInteger // Next character after position
	lineNumber   protocol.UInteger
	lineStart    protocol.UInteger
	ch           rune
}

func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readRune()
	return l
}

func (l *lexer) readRune() {
	if l.readPosition >= uint32(len(l.input)) {
		l.ch = 0 // EOF
		l.position = l.readPosition
		return
	}

	r, width := utf8.DecodeRuneInString(l.input[l.readPosition:])

	l.ch = r
	l.position = l.readPosition
	l.readPosition += uint32(width)
}

func (l *lexer) peekRune() rune {
	if l.readPosition >= uint32(len(l.input)) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.lineNumber++
			l.readRune()
			l.lineStart = l.position
			continue
		}
		l.readRune()
	}
}

func (l *lexer) nextToken() token {
	var tok token

	l.skipWhitespace()

	switch l.ch {
	case '#':
		return l.readComment()
	}

	tok.tokenType = "eof"
	return tok
}

func (l *lexer) readComment() token {
	charPos := l.position - l.lineStart

	startPos := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readRune()
	}

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    l.position - startPos,
		tokenType: string(protocol.SemanticTokenTypeComment),
	}
}

type parser struct {
	l         *lexer
	curToken  token
	peekToken token
}

func newParser(input string) parser {
	p := parser{
		l: newLexer(input),
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.nextToken()
}

func parseQuadlet(fileText string) []token {
	tokens := []token{}

	p := newParser(fileText)

	for p.curToken.tokenType != "eof" {
		p.nextToken()
	}

	return tokens
}
