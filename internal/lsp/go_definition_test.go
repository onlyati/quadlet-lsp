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
	err := os.WriteFile(path, []byte(content), 0o644)
	assert.NoError(t, err)
	return path
}

func TestFindQuadlets_MatchingFile(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "mypod.pod", "dummy content")

	loc, err := findQuadlets("pod", "mypod.pod", tmpDir, 2)
	assert.NoError(t, err)
	assert.Contains(t, string(loc.URI), "mypod.pod")
}

func TestFindQuadlets_VolumeColonSuffix(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "myvol.volume", "dummy content")

	loc, err := findQuadlets("volume", "myvol.volume:ro", tmpDir, 2)
	assert.NoError(t, err)
	assert.Contains(t, string(loc.URI), "myvol.volume")
}

func TestFindQuadlets_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "somethingelse.network", "dummy content")

	loc, err := findQuadlets("network", "notfound.network", tmpDir, 2)
	assert.NoError(t, err)
	assert.Equal(t, "", string(loc.URI))
}
