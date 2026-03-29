package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletCommentASCII(t *testing.T) {
	input := `
# First comment line
	# Second comment line
`

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("# First comment line")),
			tokenType: string(protocol.SemanticTokenTypeComment),
			text:      "# First comment line",
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("# Second comment line")),
			tokenType: string(protocol.SemanticTokenTypeComment),
			text:      "# Second comment line",
		},
	}

	tokens := []semanticToken{}
	l := newLexer(input)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	assert.Len(t, tokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tokens {
		require.Equal(t, expected[i], token, "invalid token parsed")
	}
}

func Test_parseQuadletCommentUTF16(t *testing.T) {
	input := `
# Second comment line 🫠 emoji
# 日本語 comment
`

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("# Second comment line 🫠 emoji")),
			tokenType: string(protocol.SemanticTokenTypeComment),
			text:      "# Second comment line 🫠 emoji",
		},
		{
			line:      2,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("# 日本語 comment")),
			tokenType: string(protocol.SemanticTokenTypeComment),
			text:      "# 日本語 comment",
		},
	}

	tokens := []semanticToken{}
	l := newLexer(input)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	assert.Len(t, tokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tokens {
		require.Equal(t, expected[i], token, "invalid token parsed")
	}
}

func Test_parseQuadletSection(t *testing.T) {
	input := `
[Container]
[Unit]
`

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("[Container]")),
			tokenType: string(protocol.SemanticTokenTypeNamespace),
			text:      "[Container]",
		},
		{
			line:      2,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("[Unit]")),
			tokenType: string(protocol.SemanticTokenTypeNamespace),
			text:      "[Unit]",
		},
	}

	tokens := []semanticToken{}
	l := newLexer(input)
	tok := l.nextToken()

	for tok.tokenType != "eof" {
		tokens = append(tokens, tok)
		tok = l.nextToken()
	}

	assert.Len(t, tokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tokens {
		require.Equal(t, expected[i], token, "invalid token parsed")
	}
}

func Test_parseQuadletProperty(t *testing.T) {
	input := `
Foo=bar
Bar=foobar \
  foo
`

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Foo")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Foo",
		},
		{
			line:      1,
			charPos:   3,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      1,
			charPos:   4,
			length:    protocol.UInteger(utils.Utf16Len("bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar",
		},
		{
			line:      2,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Bar")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Bar",
		},
		{
			line:      2,
			charPos:   3,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   4,
			length:    protocol.UInteger(utils.Utf16Len("foobar ")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "foobar ",
		},
		{
			line:      2,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("\\")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\\",
		},
		{
			line:      3,
			charPos:   2,
			length:    protocol.UInteger(utils.Utf16Len("foo")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "foo",
		},
	}

	tokens := []semanticToken{}
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
