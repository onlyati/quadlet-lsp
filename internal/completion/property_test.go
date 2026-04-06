package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertyCompletion(t *testing.T) {
	scenarios := []struct {
		input    string
		position parser.NodePosition
	}{
		{
			"[Container]\nDNS=",
			parser.NodePosition{LineNumber: 1, Position: 4},
		},
		{
			"[Container]\nDNS=\nNetwork=foo.network",
			parser.NodePosition{LineNumber: 1, Position: 4},
		},
	}

	for i, scenario := range scenarios {
		p := parser.NewParserFromMemory("foo.container", scenario.input)
		tokenInfo := p.Quadlet.FindToken(scenario.position)
		require.Len(t, tokenInfo.ParentNodes, 2, "failed at scenario %d", i)

		s := NewCompletion(
			[]string{},
			"foo.container",
			scenario.position.LineNumber,
			scenario.position.Position,
			p.Quadlet,
			tokenInfo,
		)
		s.commander = imageMockCommander{}
		comps := s.RunCompletion(&utils.QuadletConfig{})
		require.Greaterf(t, len(comps), 0, "did not found completions at %d", i)

		labels := []string{}
		for _, c := range comps {
			labels = append(labels, c.Label)
		}

		assert.NotContains(t, labels, "foo.volume")
		assert.Contains(t, labels, "1.1.1.1", "did not list parameters")
	}
}
