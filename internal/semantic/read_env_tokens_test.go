package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletEnv(t *testing.T) {
	input := `Environment=FOO=bar`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Environment")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Environment",
		},
		{
			line:      0,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   12,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   15,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   16,
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

func Test_parseQuadletEnvWithAposhtrophes(t *testing.T) {
	input := `Environment="FOO=foo bar"`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Environment")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Environment",
		},
		{
			line:      0,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   12,
			length:    protocol.UInteger(utils.Utf16Len("\"")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\"",
		},
		{
			line:      0,
			charPos:   13,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   17,
			length:    protocol.UInteger(utils.Utf16Len("foo bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "foo bar",
		},
		{
			line:      0,
			charPos:   24,
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

func Test_parseQuadletEnvWithQuoteMark(t *testing.T) {
	input := `Environment='FOO=foo bar'`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Environment")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Environment",
		},
		{
			line:      0,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   12,
			length:    protocol.UInteger(utils.Utf16Len("'")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "'",
		},
		{
			line:      0,
			charPos:   13,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   17,
			length:    protocol.UInteger(utils.Utf16Len("foo bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "foo bar",
		},
		{
			line:      0,
			charPos:   24,
			length:    protocol.UInteger(utils.Utf16Len("'")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "'",
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

func Test_parseQuadletEnvComplex(t *testing.T) {
	input := `Environment=FOO=BAR FOO2=BAR2 "MyVar=MyValue is=>here" 'foo=bar'`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Environment")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Environment",
		},
		{
			line:      0,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   12,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   15,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len("BAR")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "BAR",
		},
		{
			line:      0,
			charPos:   20,
			length:    protocol.UInteger(utils.Utf16Len("FOO2")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO2",
		},
		{
			line:      0,
			charPos:   24,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   25,
			length:    protocol.UInteger(utils.Utf16Len("BAR2")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "BAR2",
		},
		{
			line:      0,
			charPos:   30,
			length:    protocol.UInteger(utils.Utf16Len("\"")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\"",
		},
		{
			line:      0,
			charPos:   31,
			length:    protocol.UInteger(utils.Utf16Len("MyVar")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "MyVar",
		},
		{
			line:      0,
			charPos:   36,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   37,
			length:    protocol.UInteger(utils.Utf16Len("MyValue is=>here")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "MyValue is=>here",
		},
		{
			line:      0,
			charPos:   53,
			length:    protocol.UInteger(utils.Utf16Len("\"")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\"",
		},
		{
			line:      0,
			charPos:   55,
			length:    protocol.UInteger(utils.Utf16Len("'")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "'",
		},
		{
			line:      0,
			charPos:   56,
			length:    protocol.UInteger(utils.Utf16Len("foo")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo",
		},
		{
			line:      0,
			charPos:   59,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   60,
			length:    protocol.UInteger(utils.Utf16Len("bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar",
		},
		{
			line:      0,
			charPos:   63,
			length:    protocol.UInteger(utils.Utf16Len("'")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "'",
		},
	}

	tokens := []token{}
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

func Test_parseQuadletEnvMultiline(t *testing.T) {
	input := `
Environment= \
	"FOO=foo bar" \
	foo=bar`

	expected := []token{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Environment")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Environment",
		},
		{
			line:      1,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("\"")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\"",
		},
		{
			line:      2,
			charPos:   2,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      2,
			charPos:   5,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("foo bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "foo bar",
		},
		{
			line:      2,
			charPos:   13,
			length:    protocol.UInteger(utils.Utf16Len("\"")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\"",
		},
		{
			line:      3,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("foo")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo",
		},
		{
			line:      3,
			charPos:   4,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      3,
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

func Test_parseQuadletEnvWithNum(t *testing.T) {
	input := `Environment=FOO=127.0.0.1`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Environment")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Environment",
		},
		{
			line:      0,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   12,
			length:    protocol.UInteger(utils.Utf16Len("FOO")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "FOO",
		},
		{
			line:      0,
			charPos:   15,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len("127.0.0.1")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "127.0.0.1",
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
