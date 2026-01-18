package completion

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

type imageMockCommander struct{}

func (c imageMockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{"image1", "image2"}, nil
}

func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0o644)
	assert.NoError(t, err)
	return path
}

func createTempDir(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.Mkdir(path, 0o755)
	assert.NoError(t, err)
	return path
}

func TestPropertyImage_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.image", "[Image]")
	createTempFile(t, tmpDir, "foo.build", "[Build]")
	createTempFile(t, tmpDir, "foo.volume", "[Volume]")

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

	if slices.Contains(labels, "foo.volume") {
		t.Fatal("cannot list images")
	}

	if !slices.Contains(labels, "image1") {
		t.Fatal("did not read commander output")
	}

	if !slices.Contains(labels, "foo.image") {
		t.Fatal("did not list image files")
	}

	if !slices.Contains(labels, "foo.build") {
		t.Fatal("did not list build files")
	}
}
