// Package lexer is a module to hold logic for Quadlet file tokenization.
package lexer

import (
	"strings"
	"unicode/utf8"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

type TokenType string

const (
	TokenTypeComment  = "comment"   // For example: '# This is a comment'
	TokenTypeSection  = "section"   // For example: '[Unit]'
	TokenTypeKeyword  = "keyword"   // Text before '='
	TokenTypeValue    = "value"     // Text after '='
	TokenTypeAssign   = "assign"    // The '=' character
	TokenTypeContSign = "cont_sign" // The '\' at the end of the line
	TokenTypeEOF      = "eof"       // End of file
)

type TokenPosition struct {
	LineNumber uint32
	Position   uint32
}

type Token struct {
	StartPos TokenPosition // Where the token start
	EndPos   TokenPosition // Where the token ends
	Position uint32        // Start index of token in the input
	Length   uint32        // Length of the token
	Type     TokenType     // What kind of token it is
	Text     string        // Content of the token
}

func newToken(
	startLineNumber, startInlinePosition uint32,
	endLineNumber, endInlinePosition uint32,
	position uint32,
	length uint32,
	tokenType TokenType,
	text string,
) Token {
	return Token{
		StartPos: TokenPosition{startLineNumber, startInlinePosition},
		EndPos:   TokenPosition{endLineNumber, endInlinePosition},
		Position: position,
		Length:   length,
		Type:     tokenType,
		Text:     text,
	}
}

type Lexer struct {
	Input             string  // Text that lexer process
	position          uint32  // Current position of the lexer
	readPosition      uint32  // Next reading position
	lineNumber        uint32  // Number of line is read in 'Position'
	lineStartPosition uint32  // Which 'Position' this line started
	character         rune    // Actual character
	Tokens            []Token // The processed tokens from the 'Input', this is the result
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		Input:             input,
		position:          0,
		readPosition:      0,
		lineNumber:        0,
		lineStartPosition: 0,
		character:         0,
		Tokens:            []Token{},
	}
	l.readRune() // Set the first character

	return l
}

// Run method run through the read file and get the tokens. Function ends when
// reach the end of the file.
func (l *Lexer) Run() {
	l.nextToken()

	for l.LastToken().Type != TokenTypeEOF {
		l.nextToken()
	}
}

func (l *Lexer) LastToken() *Token {
	if len(l.Tokens) > 0 {
		return &l.Tokens[len(l.Tokens)-1]
	}
	return nil
}

// readRune method reads the next character. UTF16 compatible.
func (l *Lexer) readRune() {
	if l.readPosition >= uint32(len(l.Input)) {
		l.character = 0 // EOF
		l.position = l.readPosition
		return
	}

	r, width := utf8.DecodeRuneInString(l.Input[l.readPosition:])

	l.character = r
	l.position = l.readPosition
	l.readPosition += uint32(width)
}

func (l *Lexer) peekRune() rune {
	if l.readPosition >= uint32(len(l.Input)) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.Input[l.readPosition:])
	return r
}

func (l *Lexer) handleNewLine() {
	l.lineNumber++
	l.readRune() // Consume '\n'
	l.lineStartPosition = l.position
}

func (l *Lexer) skipInlineWhitespace() {
	for l.character == ' ' || l.character == '\t' || l.character == '\r' {
		l.readRune()
	}
}

func (l *Lexer) skipWhitespace() {
	for l.character == ' ' || l.character == '\t' || l.character == '\n' || l.character == '\r' {
		if l.character == '\n' {
			l.handleNewLine()
			continue
		}
		l.readRune()
	}
}

type readUntilParm struct {
	delimiters map[rune]any // What character until we want to read. Default is EOF
	tokenType  TokenType    // The result token's type
	readLast   bool         // Read the character after delimiter is found
}

// readUntil is a regular reader to read until specified delimiters. This is
// used to fetch tokens like comment or section.
func (l *Lexer) readUntil(parms readUntilParm) Token {
	startByte := l.position
	startLine := l.lineNumber
	startInlinePos := utils.Utf16Len(l.Input[l.lineStartPosition:l.position])

	parms.delimiters[0] = nil // Add end of line by default for safety reasons
	for {
		_, found := parms.delimiters[l.character]
		if found {
			break
		}
		if l.character == '\n' {
			l.handleNewLine()
		} else {
			l.readRune()
		}
	}

	if parms.readLast {
		l.readRune()
	}

	endLine := l.lineNumber
	endInlinePos := utils.Utf16Len(l.Input[l.lineStartPosition:l.position])

	text := l.Input[startByte:l.position]

	return newToken(
		startLine, startInlinePos,
		endLine, endInlinePos,
		startByte,
		utils.Utf16Len(text),
		parms.tokenType,
		text,
	)
}

func (l *Lexer) readOneRune(tokenType TokenType) Token {
	startByte := l.position
	startLine := l.lineNumber
	startInlinePos := utils.Utf16Len(l.Input[l.lineStartPosition:l.position])

	l.readRune()

	endLine := l.lineNumber
	endInlinePos := utils.Utf16Len(l.Input[l.lineStartPosition:l.position])

	text := l.Input[startByte:l.position]

	return newToken(
		startLine, startInlinePos,
		endLine, endInlinePos,
		startByte,
		utils.Utf16Len(text),
		tokenType,
		text,
	)
}

func (l *Lexer) readKeyValuePair() {
	for {
		lastToken := l.LastToken()
		l.skipInlineWhitespace()

		switch l.character {
		case '\n':
			l.handleNewLine()
			if lastToken != nil {
				if lastToken.Type == TokenTypeEOF {
					return
				}
				if lastToken.Type != TokenTypeContSign {
					return
				}
			}
		case 0:
			return
		case '\\':
			l.Tokens = append(l.Tokens, l.readOneRune(TokenTypeContSign))
		case '=':
			l.Tokens = append(l.Tokens, l.readOneRune(TokenTypeAssign))
			valueToken := l.readUntil(readUntilParm{
				delimiters: map[rune]any{'\n': nil, '\\': nil},
				tokenType:  TokenTypeValue,
				readLast:   false,
			})
			if strings.TrimSpace(valueToken.Text) != "" {
				l.Tokens = append(l.Tokens, valueToken)
			}
		default:
			// If the last token was a continue sign, then read value, else keyword
			if lastToken != nil {
				if l.LastToken().Type == TokenTypeContSign {
					l.Tokens = append(l.Tokens, l.readUntil(readUntilParm{
						delimiters: map[rune]any{'\n': nil, '\\': nil},
						tokenType:  TokenTypeValue,
						readLast:   false,
					}))
					continue
				}
			}
			l.Tokens = append(l.Tokens, l.readUntil(readUntilParm{
				delimiters: map[rune]any{'=': nil, '\\': nil, '\n': nil},
				tokenType:  TokenTypeKeyword,
				readLast:   false,
			}))
		}
	}
}

func (l *Lexer) nextToken() {
	for {
		l.skipWhitespace()

		switch l.character {
		case '#':
			l.Tokens = append(l.Tokens, l.readUntil(readUntilParm{
				delimiters: map[rune]any{'\n': nil},
				tokenType:  TokenTypeComment,
				readLast:   false,
			}))
		case '[':
			l.Tokens = append(l.Tokens, l.readUntil(readUntilParm{
				delimiters: map[rune]any{']': nil},
				tokenType:  TokenTypeSection,
				readLast:   true,
			}))
		case 0:
			l.Tokens = append(l.Tokens, newToken(
				l.lineNumber,
				l.position-l.lineStartPosition,
				l.lineNumber,
				l.position-l.lineStartPosition,
				l.position,
				0,
				TokenTypeEOF,
				"",
			))
			return
		default:
			if utils.IsLetter(l.character) {
				// Here we looking for key-value pairs like 'Pod=foo.pod'
				// The keyword itself must start with a normal UTF8 character
				l.readKeyValuePair()
			} else {
				// This is something unkown, just iterate it over. Normally this
				// branch should not hit, but it is safe to be here to avoid
				// inifinite loops.
				l.readRune()
			}
		}
	}
}
