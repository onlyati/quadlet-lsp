package syntax

import (
	"testing"
)

func TestQSR001_WithValidSection(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nName=app")

	diags := qsr001(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}

func TestQSR001_WithoutValidSection(t *testing.T) {
	s := NewSyntaxChecker("Name=app\nExec=run.sh")

	diags := qsr001(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	diag := diags[0]

	if diag.Message == "" || diag.Source == nil || *diag.Source != "quadlet-lsp.qsr001" {
		t.Errorf("Unexpected diagnostic content: %+v", diag)
	}
}
