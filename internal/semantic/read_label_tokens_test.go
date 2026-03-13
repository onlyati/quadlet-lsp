package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletLabel(t *testing.T) {
	input := `Label=FOO=bar`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Label")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Label",
		},
		{
			line:      0,
			charPos:   5,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   9,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   10,
			length:    protocol.UInteger(utils.Utf16Len("bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar",
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

func Test_parseQuadletLabelMultiline(t *testing.T) {
	input := `
Label= \
	FOO=bar`

	expected := []token{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Label")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Label",
		},
		{
			line:      1,
			charPos:   5,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      2,
			charPos:   4,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   5,
			length:    protocol.UInteger(utils.Utf16Len("bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar",
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

func Test_parseQuadletLabelWithAposhtrophes(t *testing.T) {
	input := `Label="FOO=foo bar"`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Label")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Label",
		},
		{
			line:      0,
			charPos:   5,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("\"")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\"",
		},
		{
			line:      0,
			charPos:   7,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   10,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("foo bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "foo bar",
		},
		{
			line:      0,
			charPos:   18,
			length:    protocol.UInteger(utils.Utf16Len("\"")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\"",
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

func Test_parseQuadletLabelMultiValue(t *testing.T) {
	input := `Label=FOO=bar foo=bar`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Label")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Label",
		},
		{
			line:      0,
			charPos:   5,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   9,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   10,
			length:    protocol.UInteger(utils.Utf16Len("bar foo=bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar foo=bar",
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
