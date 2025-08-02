package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR003_ValidProperties(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nContainerName=app\nExec=run.sh\n# Something=comment\nUser=root",
		"test.container",
	)
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 5, 2),
	}
	diags := qsr003(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}

func TestQSR003_InvalidProperty(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nContainerName=app\nFoobar=yes\nExec=run.sh", "test.container")
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 5, 2),
	}
	diags := qsr003(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	diag := diags[0]

	if diag.Message == "" || diag.Source == nil || *diag.Source != "quadlet-lsp.qsr003" {
		t.Errorf("Unexpected diagnostic: %+v", diag)
	}

	expectedMessage := "Invalid property is found: Container.Foobar"
	if diag.Message != expectedMessage {
		t.Errorf("Unexpected message:\n  got: %s\n want: %s", diag.Message, expectedMessage)
	}
}

func TestQSR003_UnknownSection(t *testing.T) {
	s := NewSyntaxChecker("[Unit]\nDescription=42", "test.container")
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 5, 2),
	}
	diags := qsr003(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics for unknown section, got %d", len(diags))
	}
}

func TestQSR003_OldVersion(t *testing.T) {
	// Memory for container is available from 5.5.0
	s := NewSyntaxChecker("[Container]\nContainerName=app\nMemory=512M", "test.container")
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 4, 2),
	}
	diags := qsr003(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	diag := diags[0]

	if diag.Message == "" || diag.Source == nil || *diag.Source != "quadlet-lsp.qsr003" {
		t.Errorf("Unexpected diagnostic: %+v", diag)
	}

	expectedMessage := "Invalid property is found: Container.Memory"
	if diag.Message != expectedMessage {
		t.Errorf("Unexpected message:\n  got: %s\n want: %s", diag.Message, expectedMessage)
	}
}
