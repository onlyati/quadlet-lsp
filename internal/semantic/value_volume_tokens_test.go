package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletVolume(t *testing.T) {
	input := `Volume=foo.volume:/app:ro,U`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Volume")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
		},
		{
			line:      0,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      0,
			charPos:   7,
			length:    protocol.UInteger(utils.Utf16Len("foo.volume")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
		},
		{
			line:      0,
			charPos:   17,
			length:    protocol.UInteger(utils.Utf16Len(":")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      0,
			charPos:   18,
			length:    protocol.UInteger(utils.Utf16Len("/app")),
			tokenType: string(protocol.SemanticTokenTypeString),
		},
		{
			line:      0,
			charPos:   22,
			length:    protocol.UInteger(utils.Utf16Len(":")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      0,
			charPos:   23,
			length:    protocol.UInteger(utils.Utf16Len("ro")),
			tokenType: string(protocol.SemanticTokenTypeString),
		},
		{
			line:      0,
			charPos:   25,
			length:    protocol.UInteger(utils.Utf16Len(",")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      0,
			charPos:   26,
			length:    protocol.UInteger(utils.Utf16Len("U")),
			tokenType: string(protocol.SemanticTokenTypeString),
		},
	}

	tokens := []token{}
	l := newLexer(input)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	assert.Len(t, tokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}

func Test_parseQuadletVolumeMultiline(t *testing.T) {
	input := `
Volume= \
	foo.volume:/app:ro,U`

	expected := []token{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Volume")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
		},
		{
			line:      1,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("foo.volume")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
		},
		{
			line:      2,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len(":")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      2,
			charPos:   12,
			length:    protocol.UInteger(utils.Utf16Len("/app")),
			tokenType: string(protocol.SemanticTokenTypeString),
		},
		{
			line:      2,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len(":")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      2,
			charPos:   17,
			length:    protocol.UInteger(utils.Utf16Len("ro")),
			tokenType: string(protocol.SemanticTokenTypeString),
		},
		{
			line:      2,
			charPos:   19,
			length:    protocol.UInteger(utils.Utf16Len(",")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			line:      2,
			charPos:   20,
			length:    protocol.UInteger(utils.Utf16Len("U")),
			tokenType: string(protocol.SemanticTokenTypeString),
		},
	}

	tokens := []token{}
	l := newLexer(input)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	assert.Len(t, tokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}
