package semantic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadletCommentASCII(t *testing.T) {
	input := `
# Fist comment line
	# Second comment line
`

	expected := []token{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(len("# Fist comment line")),
			tokenType: string(protocol.SemanticTokenTypeComment),
		},
		{
			line:      2,
			charPos:   1,
			length:    protocol.UInteger(len("# Second comment line")),
			tokenType: string(protocol.SemanticTokenTypeComment),
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
		require.Equal(t, expected[i], token, "invalid token parsed")
	}
}

func Test_parseQuadletCommentUTF16(t *testing.T) {
	input := `
# Second comment line 🫠 emoji
# 日本語 comment
`

	expected := []token{
		{
			line:      1,
			charPos:   0,
			length:    protocol.UInteger(len("# Second comment line 🫠 emoji")),
			tokenType: string(protocol.SemanticTokenTypeComment),
		},
		{
			line:      2,
			charPos:   0,
			length:    protocol.UInteger(len("# 日本語 comment")),
			tokenType: string(protocol.SemanticTokenTypeComment),
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
		require.Equal(t, expected[i], token, "invalid token parsed")
	}
}
