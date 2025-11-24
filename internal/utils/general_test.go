package utils_test

import (
	"os"
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

func TestFirstCharacterToUpper(t *testing.T) {
	v := utils.FirstCharacterToUpper("fooBar")

	if v != "FooBar" {
		t.Fatalf("Expected 'FooBar', instead got %s", v)
	}
}

func TestListQuadletFiles(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.pod", "placeholder")
	createTempFile(t, tmpDir, "foo.network", "placeholder")
	createTempFile(t, tmpDir, "bar.pod", "placeholder")

	createTempDir(t, tmpDir, "foobar")
	createTempFile(
		t,
		tmpDir+string(os.PathSeparator)+"foobar",
		"foobar.pod",
		"placeholder",
	)

	createTempDir(
		t,
		tmpDir+string(os.PathSeparator)+"foobar",
		"foo",
	)
	createTempFile(
		t,
		tmpDir+string(os.PathSeparator)+"foobar"+string(os.PathSeparator)+"foo",
		"foo.pod",
		"placeholder",
	)

	items, err := utils.ListQuadletFiles("pod", tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(items))
}

func TestTemplateNameConversion(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "web@8080.volume",
			expected: "web@.volume",
		},
		{
			input:    "web@siteA.container",
			expected: "web@.container",
		},
	}

	for _, s := range cases {
		result := utils.ConvertTemplateNameToFile(s.input)
		if s.expected != result {
			t.Fatalf("expected '%s' but got '%s'", s.expected, result)
		}
	}
}
