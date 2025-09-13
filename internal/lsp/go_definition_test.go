package lsp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
	return path
}

func createTempDir(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.Mkdir(path, 0755)
	assert.NoError(t, err)
	return path
}

func TestFindQuadlets_MatchingFile(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "mypod.pod", "dummy content")

	loc, err := findQuadlets("*.pod", "mypod.pod")
	assert.NoError(t, err)
	assert.Contains(t, string(loc.URI), "mypod.pod")
}

func TestFindQuadlets_VolumeColonSuffix(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "myvol.volume", "dummy content")

	loc, err := findQuadlets("*.volume", "myvol.volume:ro")
	assert.NoError(t, err)
	assert.Contains(t, string(loc.URI), "myvol.volume")
}

func TestFindQuadlets_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "somethingelse.network", "dummy content")

	loc, err := findQuadlets("*.network", "notfound.network")
	assert.NoError(t, err)
	assert.Equal(t, "", string(loc.URI))
}
