package semantic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Test_parseQuadlet(t *testing.T) {
	input := `# This is a token
[Container]
Exec=tail -f /dev/null
`

	expected := []Token{
		{
			Line:      0,
			CharPos:   0,
			Length:    uint32(len("# This is a token")),
			TokenType: string(protocol.SemanticTokenTypeComment),
		},
		{
			Line:      1,
			CharPos:   0,
			Length:    uint32(len("[Container]")),
			TokenType: string(protocol.SemanticTokenTypeNamespace),
		},
		{
			Line:      2,
			CharPos:   0,
			Length:    uint32(len("Exec")),
			TokenType: string(protocol.SemanticTokenTypeProperty),
		},
		{
			Line:      2,
			CharPos:   uint32(len("Exec")),
			Length:    1,
			TokenType: string(protocol.SemanticTokenTypeOperator),
		},
		{
			Line:      2,
			CharPos:   uint32(len("Exec=")),
			Length:    uint32(len("tail -f /dev/null")),
			TokenType: string(protocol.SemanticTokenTypeString),
		},
	}

	result := parseQuadlet(input)

	assert.Equal(t, expected, result, "wrong parsing")
}
