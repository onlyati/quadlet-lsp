package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

type networkMockCommnander struct{}

func (c networkMockCommnander) Run(name string, args ...string) ([]string, error) {
	return []string{"network1", "network2"}, nil
}

// TestPropertyNetwork_ListNetwork tests if only *.network and existing networks
// are displayed in completion for Network.
func TestPropertyNetwork_ListNetwork(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.network", "[Network]")
	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]")

	s := NewCompletion(
		[]string{"Network="},
		"test.container",
		0,
		uint32(len("Network=")),
	)
	s.commander = networkMockCommnander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListNetworks(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	assert.NotContains(t, labels, "foo.volume", "listed volume but it should not")
	assert.ElementsMatch(t,
		labels,
		[]string{"network1", "network2", "foo.network"},
		"did not list everything",
	)
}
