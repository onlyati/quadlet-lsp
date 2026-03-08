package semantic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_ImagevalueTokens(t *testing.T) {
	input := "docker.io/gitea/gitea:latest-rootless"

	inputSlice := []string{"docker.io", "gitea", "gitea:latest-rootless"}
	expected := []Token{
		{
			CharPos:   0,
			Length:    uint32(len(inputSlice[0])),
			TokenType: string(protocol.SemanticTokenTypeString),
		},
		{
			CharPos:   uint32(len("docker.io")),
			Length:    1,
			TokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			CharPos:   uint32(len("docker.io/")),
			Length:    uint32(len(inputSlice[1])),
			TokenType: string(protocol.SemanticTokenTypeClass),
		},
		{
			CharPos:   uint32(len("docker.io/gitea")),
			Length:    1,
			TokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			CharPos:   uint32(len("docker.io/gitea/")),
			Length:    uint32(len(inputSlice[2])),
			TokenType: string(protocol.SemanticTokenTypeClass),
		},
	}

	result := ImageValueTokens(input)
	assert.Equal(t, expected, result, "wrong image value tokens")
}
