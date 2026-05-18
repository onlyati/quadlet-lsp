package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type imageMockCommander struct{}

func (c imageMockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{"image1", "image2"}, nil
}

// TestPropertyImage_Valid tests if only *.image, *.build and pulled images are
// showed in the completion.
func TestPropertyImage_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.image", "[Image]")
	testutils.CreateTempFile(t, tmpDir, "foo.build", "[Build]")
	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]")

	scenarios := []struct {
		input    string
		position parser.NodePosition
	}{
		{
			"[Container]\nImage=",
			parser.NodePosition{LineNumber: 1, Position: 6},
		},
	}

	for i, scenario := range scenarios {
		p := parser.NewParserFromMemory("foo.container", scenario.input)
		tokenInfo := p.Quadlet.FindToken(scenario.position)

		s := NewCompletion(
			[]string{},
			"foo.container",
			scenario.position.LineNumber,
			scenario.position.Position,
			p.Quadlet,
			tokenInfo,
		)
		s.commander = imageMockCommander{}
		comps := s.RunCompletion(&utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.AsPtr(2),
			},
		})
		require.Greaterf(t, len(comps), 0, "did not found completions at %d", i)

		labels := []string{}
		for _, c := range comps {
			labels = append(labels, c.Label)
		}

		assert.NotContains(t, labels, "foo.volume")
		assert.Contains(t, labels, "image1", "did not read commander output")
		assert.Contains(t, labels, "foo.image", "did not list image files")
		assert.Contains(t, labels, "foo.build", "did not list build files")
	}
}
