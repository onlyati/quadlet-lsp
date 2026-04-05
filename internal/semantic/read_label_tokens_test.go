package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletLabel(t *testing.T) {
	input := `Label=FOO=bar`

	expected := []semanticToken{
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
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar",
		},
	}

	parser := parser.NewParserFromMemory("foo.container", input)
	tc := tokenConverter{
		lexerTokens:    parser.LexerTokens,
		index:          -1,
		semanticTokens: []semanticToken{},
	}
	tc.parseQuadlet()

	assert.Len(t, tc.semanticTokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tc.semanticTokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}

func Test_parseQuadletLabelMultiline(t *testing.T) {
	input := `
Label= \
	FOO=bar`

	expected := []semanticToken{
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
			line:      1,
			charPos:   7,
			length:    protocol.UInteger(utils.Utf16Len("\\")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\\",
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
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar",
		},
	}

	parser := parser.NewParserFromMemory("foo.container", input)
	tc := tokenConverter{
		lexerTokens:    parser.LexerTokens,
		index:          -1,
		semanticTokens: []semanticToken{},
	}
	tc.parseQuadlet()

	assert.Len(t, tc.semanticTokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tc.semanticTokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}

func Test_parseQuadletLabelWithAposhtrophes(t *testing.T) {
	input := `Label="FOO=foo bar"`

	expected := []semanticToken{
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
			tokenType: string(protocol.SemanticTokenTypeString),
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
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "\"",
		},
	}

	parser := parser.NewParserFromMemory("foo.container", input)
	tc := tokenConverter{
		lexerTokens:    parser.LexerTokens,
		index:          -1,
		semanticTokens: []semanticToken{},
	}
	tc.parseQuadlet()

	assert.Len(t, tc.semanticTokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tc.semanticTokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}

func Test_parseQuadletLabelMultiValue(t *testing.T) {
	input := `Label=FOO=bar foo=bar`

	expected := []semanticToken{
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
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("bar foo=bar")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "bar foo=bar",
		},
	}

	parser := parser.NewParserFromMemory("foo.container", input)
	tc := tokenConverter{
		lexerTokens:    parser.LexerTokens,
		index:          -1,
		semanticTokens: []semanticToken{},
	}
	tc.parseQuadlet()

	assert.Len(t, tc.semanticTokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tc.semanticTokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}
