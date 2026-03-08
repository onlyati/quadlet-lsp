// Package semantic contains functions, structs and constants that is related
// for semantic tokens.
package semantic

import (
	"slices"
	"strings"

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
		deltaLine := token.Line - lastLine
		deltaStart := token.CharPos

		if deltaLine == 0 {
			deltaStart = token.CharPos - lastChar
		}

		data = append(data,
			deltaLine,
			deltaStart,
			token.Length,
			uint32(slices.Index(TokenLegends, token.TokenType)),
			0, // No modifiers
		)

		lastLine = token.Line
		lastChar = token.CharPos
	}

	return &protocol.SemanticTokens{Data: data}, nil
}

type Token struct {
	Line      protocol.UInteger
	CharPos   protocol.UInteger
	Length    protocol.UInteger
	TokenType string
}

func parseQuadlet(fileText string) []Token {
	tokens := []Token{}

	specialKeywords := map[string]func(string) []Token{
		"Image": ImageValueTokens,
	}

	for i, line := range strings.Split(fileText, "\n") {
		// Do nothing if line is empty
		if line == "" {
			continue
		}

		// This is a comment line
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			tokens = append(tokens, Token{
				Line:      uint32(i),
				CharPos:   uint32(0),
				Length:    uint32(len(line)),
				TokenType: string(protocol.SemanticTokenTypeComment),
			})
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			tokens = append(tokens, Token{
				Line:      uint32(i),
				CharPos:   uint32(0),
				Length:    uint32(len(line)),
				TokenType: string(protocol.SemanticTokenTypeNamespace),
			})
			continue
		}

		// Key value pair
		tmp := strings.SplitN(line, "=", 2)
		// No '=' sign, assume it is a continuation of previous line
		if len(tmp) == 1 {
			tokens = append(tokens, Token{
				Line:      uint32(i),
				CharPos:   0,
				Length:    uint32(len(line)),
				TokenType: string(protocol.SemanticTokenTypeString),
			})
			continue
		}

		// This is a normal key-value pair
		tokens = append(tokens, Token{
			Line:      uint32(i),
			CharPos:   uint32(0),
			Length:    uint32(len(tmp[0])),
			TokenType: string(protocol.SemanticTokenTypeProperty),
		})
		opPos := uint32(strings.Index(line, "="))
		tokens = append(tokens, Token{
			Line:      uint32(i),
			CharPos:   opPos,
			Length:    uint32(1),
			TokenType: string(protocol.SemanticTokenTypeOperator),
		})
		if fn, ok := specialKeywords[tmp[0]]; ok {
			valueTokens := fn(tmp[1])
			for _, token := range valueTokens {
				token.CharPos += opPos + 1
				token.Line = uint32(i)
				tokens = append(tokens, token)
			}
		} else {
			tokens = append(tokens, Token{
				Line:      uint32(i),
				CharPos:   uint32(strings.Index(line, "=") + 1),
				Length:    uint32(len(tmp[1])),
				TokenType: string(protocol.SemanticTokenTypeString),
			})
		}
		continue
	}

	return tokens
}
