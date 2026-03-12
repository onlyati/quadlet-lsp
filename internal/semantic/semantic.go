// Package semantic contains functions, structs and constants that is related
// for semantic tokens.
package semantic

import (
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
	"Image":   (*lexer).readImageValue,
	"Volume":  (*lexer).readVolumeValue,
	"Pod":     (*lexer).readPodValue,
	"Network": (*lexer).readNetworkValue,
	"Secret":  (*lexer).readSecretValue,
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
