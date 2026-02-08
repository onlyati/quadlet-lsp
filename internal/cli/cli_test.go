package cli

import (
	"os"
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRunCLI(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	testutils.CreateTempFile(t, tmpDir, ".quadletrc.json", `{ "podmanVersion": "5.5.0" }`)
	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nImage=foo")
	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.container.d"), "label.conf", "[Container]\nLabel=Test")

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
