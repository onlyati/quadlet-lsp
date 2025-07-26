package syntax

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
	return path
}

func TestQSR006_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Container]\nImage=foo.image", "foo.container")
	diags := qsr006(s)

	if len(diags) != 0 {
		t.Fatalf("Expected no diagnostics, but got %d", len(diags))
	}
}

func TestQSR006_ValidVolume(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Volume]\nImage=foo.image", "foo.volume")
	diags := qsr006(s)

	if len(diags) != 0 {
		t.Fatalf("Expected no diagnostics, but got %d", len(diags))
	}
}

func TestQSR006_Skipped(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	s := NewSyntaxChecker("[Container]\nImage=library/debian", "foo.container")
	diags := qsr006(s)

	if len(diags) != 0 {
		t.Fatalf("Expected no diagnostics, but got %d", len(diags))
	}
}

func TestQSR006_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTestFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Container]\nImage=bar.image", "foo.container")
	diags := qsr006(s)

	if len(diags) != 1 {
		t.Fatalf("Expected no diagnostics, but got %d", len(diags))
	}

	msg := "Image file does not exists: bar.image"
	if diags[0].Message != msg {
		t.Fatalf("Wrong error message expected: '%s', got: '%s'", msg, diags[0].Message)
	}
}
