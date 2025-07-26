package syntax

import (
	"testing"
)

func TestQSR003_ValidProperties(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nContainerName=app\nExec=run.sh\nUser=root", "test.container")
	diags := qsr003(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}

func TestQSR003_InvalidProperty(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nContainerName=app\nFoobar=yes\nExec=run.sh", "test.container")
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
	diags := qsr003(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics for unknown section, got %d", len(diags))
	}
}
