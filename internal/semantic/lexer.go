package semantic

import (
	"unicode/utf8"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type token struct {
	line      protocol.UInteger
	charPos   protocol.UInteger
	length    protocol.UInteger
	tokenType string
	text      string
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
	l := &lexer{
		input:        input,
		position:     0,
		readPosition: 0,
		lineNumber:   0,
		lineStart:    0,
		ch:           0,
		queue:        []token{},
	}
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
			return token{
				tokenType: "eof",
				line:      l.lineNumber,
				charPos:   l.position,
				length:    0,
				text:      "",
			}
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
