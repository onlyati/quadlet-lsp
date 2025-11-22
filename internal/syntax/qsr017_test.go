package syntax

import (
	"os"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR017_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"test.pod",
		"[Pod]",
	)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"file://"+tmpDir+"/test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr017(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR017_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"file://"+tmpDir+"/test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr017(s)

		if len(diags) != 1 {
			t.Fatalf("Expected 1 diagnostics, got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr017" {
			t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Pod file does not exists: test.pod" {
			t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}

// TestQSR017_ValidAdjacentFile tests that pod files are found adjacent to the container file,
// even when the working directory is different (e.g., the workspace root).
func TestQSR017_ValidAdjacentFile(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Create a subdirectory for the container and pod files
	subDir := createTempDir(t, tmpDir, "subdir")

	// Create the pod file in the subdirectory
	createTempFile(
		t,
		subDir,
		"test.pod",
		"[Pod]",
	)

	// Create the container file in the subdirectory that references the pod file
	s := NewSyntaxChecker(
		"[Container]\nPod=test.pod",
		"file://"+subDir+"/test1.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
	}

	diags := qsr017(s)

	if len(diags) != 0 {
		t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
	}
}
