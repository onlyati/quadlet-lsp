package lsp

import (
	"os"
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestRunCLI(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, ".quadletrc.json", `{ "podmanVersion": "5.5.0" }`)
	createTestFile(t, tmpDir, "foo.container", "[Container]\nImage=foo")
	createTempDir(t, tmpDir, "foo.container.d")
	createTempFile(t, path.Join(tmpDir, "foo.container.d"), "label.conf", "[Container]\nLabel=Test")

	args := [][]string{
		{"quadlet-lsp", "check"},
		{"quadlet-lsp", "check", "."},
		{"quadlet-lsp", "check", tmpDir},
	}
	for _, arg := range args {
		rc, output := runCheckCLI(arg, utils.CommandExecutor{})

		if len(output) != 3 {
			t.Fatalf("expected 2 line output but got %d at %v", len(output), arg)
		}

		if rc != 4 {
			t.Fatalf("expected to get rc=4 but got rc=%d at %v", rc, arg)
		}

	}
}
