package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_ListNewProperties tests if properties are popup during completion.
func Test_ListNewProperties(t *testing.T) {
	scenarios := []struct {
		input    string
		position parser.NodePosition
		name     string
	}{
		{
			"[Container]\n",
			parser.NodePosition{LineNumber: 1, Position: 0},
			"Simple completion",
		},
		{
			"[Container]\nImage=foo.bar\n",
			parser.NodePosition{LineNumber: 1, Position: 0},
			"Completion for keyword",
		},
		{
			"[Container]\n\nImage=foo.bar\n\n[Unit]\n",
			parser.NodePosition{LineNumber: 1, Position: 0},
			"Completion mid-file",
		},
		{
			"[Container]\n\nImage=foo.bar\n\n[Unit]\n",
			parser.NodePosition{LineNumber: 2, Position: 0},
			"Completion mid empty line",
		},
		{
			"[Container]\nImage=foo.bar\n",
			parser.NodePosition{LineNumber: 2, Position: 0},
			"Last line completion",
		},
	}

	for i, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			p := parser.NewParserFromMemory("foo.container", scenario.input)
			require.Len(t, p.Errors, 0)

			s := NewCompletion(
				[]string{},
				"foo.container",
				scenario.position.LineNumber,
				scenario.position.Position,
				p.Quadlet,
				p.Quadlet.FindToken(scenario.position),
			)
			comps := s.RunCompletion(&utils.QuadletConfig{})
			assert.Greaterf(t, len(comps), 0, "did not found completions at %d", i)
		})
	}
}
