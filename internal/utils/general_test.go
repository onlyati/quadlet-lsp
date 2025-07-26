package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

func createTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
	return path
}

func TestFirstCharacterToUpper(t *testing.T) {
	v := utils.FirstCharacterToUpper("fooBar")

	if v != "FooBar" {
		t.Fatalf("Expected 'FooBar', instead got %s", v)
	}
}

func TestListQuadletFiles(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "foo.pod", "placeholder")
	createTestFile(t, tmpDir, "foo.network", "placeholder")
	createTestFile(t, tmpDir, "bar.pod", "placeholder")

	items, err := utils.ListQuadletFiles("*.pod")
	assert.NoError(t, err)

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, but got %d", len(items))
	}
}
