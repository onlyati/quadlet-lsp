// Package semantic contains functions, structs and constants that is related
// for semantic tokens.
package semantic

import (
	"strings"

	quadlet_lexer "github.com/onlyati/quadlet-lsp/pkg/quadlet/lexer"
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
	"Image":       (*lexer).readImageValue,
	"Volume":      (*lexer).readVolumeValue,
	"Pod":         (*lexer).readPodValue,
	"Network":     (*lexer).readNetworkValue,
	"Secret":      (*lexer).readSecretValue,
	"Environment": (*lexer).readEnvValue,
	"Label":       (*lexer).readLabelValue,
	"Annotation":  (*lexer).readLabelValue,
}

func CalculateSemanticTokens(lexerTokens []quadlet_lexer.Token) (*protocol.SemanticTokens, error) {
	converter := tokenConverter{
		lexerTokens:    lexerTokens,
		index:          -1,
		semanticTokens: make([]semanticToken, 0, len(lexerTokens)),
	}
	converter.parseQuadlet()

	data := make([]uint32, 0, len(converter.semanticTokens)*5)

	var lastLine, lastChar uint32

	for _, token := range converter.semanticTokens {
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

type tokenConverter struct {
	lexerTokens    []quadlet_lexer.Token
	index          int
	semanticTokens []semanticToken
}

var specialParsers = map[string]func(*tokenConverter, *quadlet_lexer.Token) protocol.SemanticTokenType{
	"Network": (*tokenConverter).readNetworkValue,
}

func (t *tokenConverter) readToken() *quadlet_lexer.Token {
	t.index++
	if t.index == len(t.lexerTokens)-1 {
		return nil
	}
	return &t.lexerTokens[t.index]
}

// parseQuadlet translate the regular lexer tokens to semantic tokens
func (t *tokenConverter) parseQuadlet() {
	token := t.readToken()
	lastKeyword := ""

	for token != nil {
		switch token.Type {
		case quadlet_lexer.TokenTypeComment:
			comment := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeComment),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, comment)
		case quadlet_lexer.TokenTypeSection:
			section := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeNamespace),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, section)
		case quadlet_lexer.TokenTypeKeyword:
			keyword := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeKeyword),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, keyword)
			lastKeyword = token.Text
		case quadlet_lexer.TokenTypeValue:
			tokenType := protocol.SemanticTokenTypeString
			if fn, ok := specialParsers[lastKeyword]; ok {
				tokenType = fn(t, token)
			}
			value := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(tokenType),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, value)
		case quadlet_lexer.TokenTypeAssign, quadlet_lexer.TokenTypeContSign:
			op := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, op)
		}

		token = t.readToken()
	}
}

// readNetworkValue is part of semantic token parsing for Network keyword.
func (t *tokenConverter) readNetworkValue(lexerToken *quadlet_lexer.Token) protocol.SemanticTokenType {
	if strings.HasSuffix(lexerToken.Text, ".network") {
		return protocol.SemanticTokenTypeParameter
	}
	return protocol.SemanticTokenTypeString
}
