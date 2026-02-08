package syntax

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0o644)
	assert.NoError(t, err)
	return path
}

func TestQSR006_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	testutils.CreateTempFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Container]\nImage=foo.image", "foo.container")
	diags := qsr006(s)
	require.Len(t, diags, 0)
}

func TestQSR006_ValidVolume(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	testutils.CreateTempFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Volume]\nImage=foo.image", "foo.volume")
	diags := qsr006(s)
	require.Len(t, diags, 0)
}

func TestQSR006_Skipped(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	s := NewSyntaxChecker("[Container]\nImage=library/debian", "foo.container")
	diags := qsr006(s)
	require.Len(t, diags, 0)
}

func TestQSR006_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Container]\nImage=bar.image", "foo.container")
	diags := qsr006(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr006", *diags[0].Source)
	assert.Equal(t, "Image file does not exists: bar.image", diags[0].Message)
}
