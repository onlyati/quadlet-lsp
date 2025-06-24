package lsp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
	return path
}

func TestFindReferences(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "example.volume", "Volume=example.volume\\nAnotherLine")

	locations, err := findReferences("Volume", "example.volume")
	assert.NoError(t, err)
	assert.Len(t, locations, 1)
	assert.Contains(t, string(locations[0].URI), "example.volume")
	assert.Equal(t, uint32(0), locations[0].Range.Start.Line)
}
