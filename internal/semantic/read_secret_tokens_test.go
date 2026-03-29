package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletSecret(t *testing.T) {
	input := `Secret=secret1,type=env,target=ENV1`

	expected := []semanticToken{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Secret")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Secret",
		},
		{
			line:      0,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   7,
			length:    protocol.UInteger(utils.Utf16Len("secret1")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "secret1",
		},
		{
			line:      0,
			charPos:   14,
			length:    protocol.UInteger(utils.Utf16Len(",")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ",",
		},
		{
			line:      0,
			charPos:   15,
			length:    protocol.UInteger(utils.Utf16Len("type")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "type",
		},
		{
			line:      0,
			charPos:   19,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   20,
			length:    protocol.UInteger(utils.Utf16Len("env")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "env",
		},
		{
			line:      0,
			charPos:   23,
			length:    protocol.UInteger(utils.Utf16Len(",")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ",",
		},
		{
			line:      0,
			charPos:   24,
			length:    protocol.UInteger(utils.Utf16Len("target")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "target",
		},
		{
			line:      0,
			charPos:   30,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   31,
			length:    protocol.UInteger(utils.Utf16Len("ENV1")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "ENV1",
		},
	}

	tokens := []semanticToken{}
	l := newLexer(input)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	require.Len(t, tokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}

func Test_parseQuadletSecretMultiline(t *testing.T) {
	input := `
Secret= \
	secret1,type=env,target=ENV1`

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Secret")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Secret",
		},
		{
			line:      1,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("secret1")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "secret1",
		},
		{
			line:      2,
			charPos:   8,
			length:    protocol.UInteger(utils.Utf16Len(",")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ",",
		},
		{
			line:      2,
			charPos:   9,
			length:    protocol.UInteger(utils.Utf16Len("type")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "type",
		},
		{
			line:      2,
			charPos:   13,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   14,
			length:    protocol.UInteger(utils.Utf16Len("env")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "env",
		},
		{
			line:      2,
			charPos:   17,
			length:    protocol.UInteger(utils.Utf16Len(",")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ",",
		},
		{
			line:      2,
			charPos:   18,
			length:    protocol.UInteger(utils.Utf16Len("target")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "target",
		},
		{
			line:      2,
			charPos:   24,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   25,
			length:    protocol.UInteger(utils.Utf16Len("ENV1")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "ENV1",
		},
	}

	tokens := []semanticToken{}
	l := newLexer(input)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	require.Len(t, tokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}
