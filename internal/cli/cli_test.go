package cli

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

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

func TestRunCLI(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, ".quadletrc.json", `{ "podmanVersion": "5.5.0" }`)
	createTempFile(t, tmpDir, "foo.container", "[Container]\nImage=foo")
	createTempDir(t, tmpDir, "foo.container.d")
	createTempFile(t, path.Join(tmpDir, "foo.container.d"), "label.conf", "[Container]\nLabel=Test")

	args := [][]string{
		{},
		{"."},
		{tmpDir},
	}
	for _, arg := range args {
		output, err := runCheckCLI(arg, utils.CommandExecutor{})

		assert.Equal(t, 3, len(output))
		assert.NoError(t, err)
	}
}
