package semantic

import (
	"testing"

	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_ImagevalueTokens(t *testing.T) {
	input := "docker.io/gitea/gitea:latest-rootless"

	expected := []token{
		{
			charPos:   0,
			length:    uint32(len("docker.io")),
			tokenType: string(protocol.SemanticTokenTypeString),
		},
		{
			charPos:   uint32(len("docker.io")),
			length:    1,
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			charPos:   uint32(len("docker.io/")),
			length:    uint32(len("gitea")),
			tokenType: string(protocol.SemanticTokenTypeClass),
		},
		{
			charPos:   uint32(len("docker.io/gitea")),
			length:    1,
			tokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			charPos:   uint32(len("docker.io/gitea/")),
			length:    uint32(len("gitea:latest-rootless")),
			tokenType: string(protocol.SemanticTokenTypeClass),
		},
	}

	result := ImageValueTokens(input)

	for i, r := range result {
		require.Equal(t, expected[i], r, "invalid token at %d", i)
	}
}
