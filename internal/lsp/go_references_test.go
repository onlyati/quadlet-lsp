package lsp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0o644)
	assert.NoError(t, err)
	return path
}

func TestFindReferences(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "example.container", "[Container]\nNetwork=example.network\nAnotherLine")

	locations, err := findReferences(
		goReferenceProperty{
			property: "Network",
			searchIn: []string{"container", "pod", "kube"},
		}, "example.network", tmpDir)
	assert.NoError(t, err)
	assert.Len(t, locations, 1)
	assert.Contains(t, string(locations[0].URI), "example.container")
	assert.Equal(t, uint32(1), locations[0].Range.Start.Line)
}

func TestFiindReferencesTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "web@.container", "[Container]\nVolume=web@%i.volume:/app")
	createTempFile(t, tmpDir, "builder@.container", "[Container]\nVolume=web@%i.volume:/app")

	locations, err := findReferences(
		goReferenceProperty{
			property: "Volume",
			searchIn: []string{"container", "pod"},
		}, "web@.volume", tmpDir)
	assert.NoError(t, err)
	assert.Len(t, locations, 2)

	for _, loc := range locations {
		if !strings.Contains(loc.URI, "web@.container") && !strings.Contains(loc.URI, "builder@.container") {
			t.Fatalf("Unexpected finding: %+v", loc)
		}
	}
}
