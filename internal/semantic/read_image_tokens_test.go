package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_converionQuadletImage(t *testing.T) {
	input := `Image=docker.io/gitea/gitea:rootless@sha256asdasdasdasd
Foo=bar`

	expected := []semanticToken{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Image")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Image",
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
			length:    protocol.UInteger(utils.Utf16Len("docker.io")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "docker.io",
		},
		{
			line:      0,
			charPos:   15,
			length:    protocol.UInteger(utils.Utf16Len("/")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "/",
		},
		{
			line:      0,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len("gitea")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "gitea",
		},
		{
			line:      0,
			charPos:   21,
			length:    protocol.UInteger(utils.Utf16Len("/")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "/",
		},
		{
			line:      0,
			charPos:   22,
			length:    protocol.UInteger(utils.Utf16Len("gitea")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "gitea",
		},
		{
			line:      0,
			charPos:   27,
			length:    protocol.UInteger(utils.Utf16Len(":")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ":",
		},
		{
			line:      0,
			charPos:   28,
			length:    protocol.UInteger(utils.Utf16Len("rootless")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "rootless",
		},
		{
			line:      0,
			charPos:   36,
			length:    protocol.UInteger(utils.Utf16Len("@")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "@",
		},
		{
			line:      0,
			charPos:   37,
			length:    protocol.UInteger(utils.Utf16Len("sha256asdasdasdasd")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "sha256asdasdasdasd",
		},
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

func Test_conversionQuadletImageValueWithoutHash(t *testing.T) {
	input := "Image=docker.io/gitea/gitea.container:rootless"

	expected := []semanticToken{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Image")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Image",
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
			length:    protocol.UInteger(utils.Utf16Len("docker.io")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "docker.io",
		},
		{
			line:      0,
			charPos:   15,
			length:    protocol.UInteger(utils.Utf16Len("/")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "/",
		},
		{
			line:      0,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len("gitea")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "gitea",
		},
		{
			line:      0,
			charPos:   21,
			length:    protocol.UInteger(utils.Utf16Len("/")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "/",
		},
		{
			line:      0,
			charPos:   22,
			length:    protocol.UInteger(utils.Utf16Len("gitea.container")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "gitea.container",
		},
		{
			line:      0,
			charPos:   37,
			length:    protocol.UInteger(utils.Utf16Len(":")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ":",
		},
		{
			line:      0,
			charPos:   38,
			length:    protocol.UInteger(utils.Utf16Len("rootless")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "rootless",
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

func Test_conversionQuadletImageValueMultiline(t *testing.T) {
	input := `
Image= \
	docker.io/gitea/gitea:rootless@sha256asdasdasdasd`

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Image")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Image",
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
			length:    protocol.UInteger(utils.Utf16Len("docker.io")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "docker.io",
		},
		{
			line:      2,
			charPos:   10,
			length:    protocol.UInteger(utils.Utf16Len("/")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "/",
		},
		{
			line:      2,
			charPos:   11,
			length:    protocol.UInteger(utils.Utf16Len("gitea")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "gitea",
		},
		{
			line:      2,
			charPos:   16,
			length:    protocol.UInteger(utils.Utf16Len("/")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "/",
		},
		{
			line:      2,
			charPos:   17,
			length:    protocol.UInteger(utils.Utf16Len("gitea")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "gitea",
		},
		{
			line:      2,
			charPos:   22,
			length:    protocol.UInteger(utils.Utf16Len(":")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ":",
		},
		{
			line:      2,
			charPos:   23,
			length:    protocol.UInteger(utils.Utf16Len("rootless")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "rootless",
		},
		{
			line:      2,
			charPos:   31,
			length:    protocol.UInteger(utils.Utf16Len("@")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "@",
		},
		{
			line:      2,
			charPos:   32,
			length:    protocol.UInteger(utils.Utf16Len("sha256asdasdasdasd")),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      "sha256asdasdasdasd",
		},
	}

	parser := parser.NewParserFromMemory("foo.container", input)
	tc := tokenConverter{
		lexerTokens:    parser.LexerTokens,
		index:          -1,
		semanticTokens: []semanticToken{},
	}
	tc.parseQuadlet()

	require.Len(t, tc.semanticTokens, len(expected), "invalid number of elements in tokens")
	for i, token := range tc.semanticTokens {
		require.Equal(t, expected[i], token, "invalid token parsed at %d.", i)
	}
}

func Test_conversionQuadletImageFile(t *testing.T) {
	input := `Image=foo.image`
	parser := parser.NewParserFromMemory("foo.container", input)

	expected := []semanticToken{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Image")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Image",
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
			length:    protocol.UInteger(utils.Utf16Len("foo.image")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo.image",
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

func Test_conversionQuadletImageFileMultiline(t *testing.T) {
	input := `
Image=\
	foo.image`
	parser := parser.NewParserFromMemory("foo.container", input)

	expected := []semanticToken{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Image")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Image",
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
			charPos:   6,
			length:    protocol.UInteger(utils.Utf16Len("\\")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "\\",
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("foo.image")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo.image",
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
