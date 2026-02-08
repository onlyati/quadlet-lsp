package lsp

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFindQuadlets_MatchingFile(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "mypod.pod", "dummy content")

	loc, err := findQuadlets("pod", "mypod.pod", tmpDir, 2)
	assert.NoError(t, err)
	assert.Contains(t, string(loc.URI), "mypod.pod")
}

func TestFindQuadlets_VolumeColonSuffix(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "myvol.volume", "dummy content")

	loc, err := findQuadlets("volume", "myvol.volume:ro", tmpDir, 2)
	assert.NoError(t, err)
	assert.Contains(t, string(loc.URI), "myvol.volume")
}

func TestFindQuadlets_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "somethingelse.network", "dummy content")

	loc, err := findQuadlets("network", "notfound.network", tmpDir, 2)
	assert.NoError(t, err)
	assert.Equal(t, "", string(loc.URI))
}
