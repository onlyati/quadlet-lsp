package utils_test

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_QuadletWalkDir(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, ".hidden")
	testutils.CreateTempDir(t, tmpDir, "dirA")
	testutils.CreateTempDir(t, path.Join(tmpDir, "dirA"), "dirAB")
	testutils.CreateTempDir(t, path.Join(tmpDir, "dirA", "dirAB"), "dirABC")

	testutils.CreateTempFile(t, tmpDir, "file.txt", "placeholder")
	testutils.CreateTempFile(t, path.Join(tmpDir, ".hidden"), "hidden.txt", "placeholder")
	testutils.CreateTempFile(t, path.Join(tmpDir, "dirA"), "fileA.txt", "placeholder")
	testutils.CreateTempFile(t, path.Join(tmpDir, "dirA", "dirAB"), "fileAB.txt", "placeholder")
	testutils.CreateTempFile(t, path.Join(tmpDir, "dirA", "dirAB", "dirABC"), "fileABC.txt", "placeholder")

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

	testutils.CreateTempDir(t, tmpDir, "foobar")
	testutils.CreateTempDir(t, tmpDir+string(os.PathSeparator)+"foobar", "foo")

	testutils.CreateTempFile(t, tmpDir, "foo.pod", "placeholder")
	testutils.CreateTempFile(t, tmpDir, "foo.network", "placeholder")
	testutils.CreateTempFile(t, tmpDir, "bar.pod", "placeholder")
	testutils.CreateTempFile(t, tmpDir+string(os.PathSeparator)+"foobar", "foobar.pod", "placeholder")
	testutils.CreateTempFile(t, tmpDir+string(os.PathSeparator)+"foobar"+string(os.PathSeparator)+"foo", "foo.pod", "placeholder")

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
