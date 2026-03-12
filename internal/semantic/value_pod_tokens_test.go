package semantic

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletPod(t *testing.T) {
	input := `Pod=foo.pod`

	expected := []token{
		{
			line:      0,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Pod")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Pod",
		},
		{
			line:      0,
			charPos:   3,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      0,
			charPos:   4,
			length:    protocol.UInteger(utils.Utf16Len("foo.pod")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo.pod",
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

func Test_parseQuadletPodMultiline(t *testing.T) {
	input := `
Pod=\
	foo.pod`

	expected := []token{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(utils.Utf16Len("Pod")),
			tokenType: string(protocol.SemanticTokenTypeKeyword),
			text:      "Pod",
		},
		{
			line:      1,
			charPos:   3,
			length:    protocol.UInteger(utils.Utf16Len("=")),
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      "=",
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(utils.Utf16Len("foo.pod")),
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      "foo.pod",
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
