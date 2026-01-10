package utils_test

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	assert.Equal(t, v, "FooBar")
}

func Test_QuadletWalkDir(t *testing.T) {
	tmpDir := t.TempDir()

	createTempFile(t, tmpDir, "file.txt", "placeholder")

	createTempDir(t, tmpDir, ".hidden")
	createTempFile(t, path.Join(tmpDir, ".hidden"), "hidden.txt", "placeholder")

	createTempDir(t, tmpDir, "dirA")
	createTempFile(t, path.Join(tmpDir, "dirA"), "fileA.txt", "placeholder")

	createTempDir(t, path.Join(tmpDir, "dirA"), "dirAB")
	createTempFile(t, path.Join(tmpDir, "dirA", "dirAB"), "fileAB.txt", "placeholder")

	createTempDir(t, path.Join(tmpDir, "dirA", "dirAB"), "dirABC")
	createTempFile(t, path.Join(tmpDir, "dirA", "dirAB", "dirABC"), "fileABC.txt", "placeholder")

	expected := []string{"file.txt", "dirA/fileA.txt", "dirA/dirAB/fileAB.txt"}
	fileList := []string{}
	err := utils.QuadletWalkDir(tmpDir, 2, func(path string, d fs.DirEntry, err error) error {
		relPath := strings.TrimPrefix(path, tmpDir)
		relPath = strings.TrimPrefix(relPath, string(os.PathSeparator))
		if d.IsDir() {
			return nil
		}
		fileList = append(fileList, relPath)
		return nil
	})

	assert.NoError(t, err)
	assert.Len(t, fileList, len(expected))
	assert.ElementsMatch(
		t,
		expected,
		fileList,
	)
}

func TestListQuadletFiles(t *testing.T) {
	tmpDir := t.TempDir()

	createTempFile(t, tmpDir, "foo.pod", "placeholder")
	createTempFile(t, tmpDir, "foo.network", "placeholder")
	createTempFile(t, tmpDir, "bar.pod", "placeholder")

	createTempDir(t, tmpDir, "foobar")
	createTempFile(t, tmpDir+string(os.PathSeparator)+"foobar", "foobar.pod", "placeholder")

	createTempDir(t, tmpDir+string(os.PathSeparator)+"foobar", "foo")
	createTempFile(t, tmpDir+string(os.PathSeparator)+"foobar"+string(os.PathSeparator)+"foo", "foo.pod", "placeholder")

	items, err := utils.ListQuadletFiles("pod", tmpDir, 2)
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
		assert.Equal(t, s.expected, result)
	}
}
