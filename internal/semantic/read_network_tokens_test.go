package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_conversionQuadletNetwork(t *testing.T) {
	input := `Network=foo.network`
	parser := parser.NewParserFromMemory("foo.container", input)

	expected := []semanticToken{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Network")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Network",
		},
		{
			line:      0,
			charPos:   7,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   8,
			length:    protocol.UInteger(utils.Utf16Len("foo.network")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo.network",
		},
	}

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

func Test_conversionQuadletNetworkMultiLine(t *testing.T) {
	input := `
Network=\
	foo.network`
	parser := parser.NewParserFromMemory("foo.container", input)

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Network")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Network",
		},
		{
			line:      1,
			charPos:   7,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      1,
			charPos:   8,
			length:    protocol.UInteger(utils.Utf16Len("\\")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\\",
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("foo.network")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo.network",
		},
	}

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
