package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

// TestPropertyPod_Valid tests if only *.pod files are display in the complation
// for Pod.
func TestPropertyPod_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.pod", "[Pod]")
	testutils.CreateTempFile(t, tmpDir, "bar.pod", "[Pod]")
	testutils.CreateTempFile(t, tmpDir, "foo.network", "[Network]")

	s := Completion{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListPods(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	assert.NotContains(t, labels, "foo.network", "listed network but should not")
	assert.ElementsMatch(
		t,
		labels,
		[]string{"foo.pod", "bar.pod"},
		"did not list everything",
	)
}
