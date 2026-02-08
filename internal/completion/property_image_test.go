package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
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

	s := Completion{}
	s.commander = imageMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListImages(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	assert.NotContains(t, labels, "foo.volume")
	assert.Contains(t, labels, "image1", "did not read commander output")
	assert.Contains(t, labels, "foo.image", "did not list image files")
	assert.Contains(t, labels, "foo.build", "did not list build files")
}
