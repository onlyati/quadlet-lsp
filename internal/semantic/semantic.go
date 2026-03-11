// Package semantic contains functions, structs and constants that is related
// for semantic tokens.
package semantic

import (
	"unicode/utf8"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var TokenLegends = []string{
	string(protocol.SemanticTokenTypeComment),   // Comment lines
	string(protocol.SemanticTokenTypeKeyword),   // Like 'Image=', 'Exec=', 'Pod='
	string(protocol.SemanticTokenTypeNamespace), // Section like '[Container]', '[Unit]'
	string(protocol.SemanticTokenTypeString),    // Value belongs to keywords
	string(protocol.SemanticTokenTypeOperator),  // Operators like '=', ':', ','
	string(protocol.SemanticTokenTypeClass),     // Used within values to highlight things
	string(protocol.SemanticTokenTypeParameter), // Used within values to highlight things
}

var LegendMap = func() map[string]uint32 {
	m := make(map[string]uint32)
	for i, t := range TokenLegends {
		m[t] = uint32(i)
	}
	return m
}()

var specialFunctionMap = map[string]func(*lexer){
	"Image":  (*lexer).readImageValue,
	"Volume": (*lexer).readVolumeValue,
	"Pod":    (*lexer).readPodValue,
}

func CalculateSemanticTokens(fileText string) (*protocol.SemanticTokens, error) {
	tokens := parseQuadlet(fileText)

	data := make([]uint32, 0, len(tokens)*5)

	var lastLine, lastChar uint32

	for _, token := range tokens {
		deltaLine := token.line - lastLine
		deltaStart := token.charPos
		if deltaLine == 0 {
			deltaStart = token.charPos - lastChar
		}

		typeIndex, ok := LegendMap[token.tokenType]
		if !ok {
			typeIndex = 0
		}

		data = append(data,
			deltaLine,
			deltaStart,
			token.length,
			typeIndex,
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
	queue        []token
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

// func (l *lexer) peekRune() rune {
// 	if l.readPosition >= uint32(len(l.input)) {
// 		return 0
// 	}
// 	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
// 	return r
// }

func (l *lexer) handleNewLine() {
	l.lineNumber++
	l.readRune()
	l.lineStart = l.position
}

func (l *lexer) skipInlineWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readRune()
	}
}

func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.handleNewLine()
			continue
		}
		l.readRune()
	}
}

func (l *lexer) nextToken() token {
	// if something has been put into the queue, then empty it
	if len(l.queue) > 0 {
		tok := l.queue[0]
		l.queue = l.queue[1:]
		return tok
	}

	for {
		l.skipWhitespace()

		switch l.ch {
		case '#':
			return l.readComment()
		case '[':
			return l.readSection()
		case '\\':
			return l.readOperator()
		case 0:
			return token{tokenType: "eof"}
		default:
			if utils.IsLetter(l.ch) {
				l.readAssignment()
				return l.nextToken()
			}

			// This is something unknow, just iterate over it until we find
			// something interesting
			l.readRune()
			continue
		}
	}
}

func (l *lexer) readOperator() token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	l.readRune()

	// Measure the content we just read in UTF-16 units
	sectionText := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(sectionText),
		tokenType: string(protocol.SemanticTokenTypeOperator),
	}
}

func (l *lexer) readValue() token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	for l.ch != '\n' && l.ch != '\\' && l.ch != 0 {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	sectionText := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(sectionText),
		tokenType: string(protocol.SemanticTokenTypeString),
	}
}

func (l *lexer) readAssignment() {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	for l.ch != '=' && l.ch != '\n' && l.ch != '\\' && l.ch != 0 {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	propText := l.input[startByte:l.position]

	// If we hit a \n, it means it was a value line, if we hit '=' it is property
	if l.ch == '\n' || l.ch == '\\' || l.ch == 0 {
		l.queue = append(l.queue, token{
			line:      l.lineNumber,
			charPos:   charPos,
			length:    utils.Utf16Len(propText),
			tokenType: string(protocol.SemanticTokenTypeString),
		})
		return
	}

	// We hit '=' a sign, put property to queue, then analyze line further
	l.queue = append(l.queue, token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(propText),
		tokenType: string(protocol.SemanticTokenTypeKeyword),
	})

	if l.ch == '=' {
		l.queue = append(l.queue, l.readOperator())

		// Check if we had any special parse for property, if not just read value
		if fn, ok := specialFunctionMap[propText]; ok {
			fn(l)
		} else {
			l.queue = append(l.queue, l.readValue())
		}
	}
}

func (l *lexer) readSection() token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	for l.ch != ']' && l.ch != '\n' && l.ch != 0 {
		l.readRune()
	}
	if l.ch == ']' {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	sectionText := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(sectionText),
		tokenType: string(protocol.SemanticTokenTypeNamespace),
	}
}

func (l *lexer) readComment() token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	for l.ch != '\n' && l.ch != 0 {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	commentText := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(commentText),
		tokenType: string(protocol.SemanticTokenTypeComment),
	}
}

func parseQuadlet(fileText string) []token {
	tokens := []token{}

	l := newLexer(fileText)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	return tokens
}
