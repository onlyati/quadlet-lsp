package syntax

import (
	"testing"
)

func TestQSR002_UnfinishedLine(t *testing.T) {
	s := NewSyntaxChecker("Name=\nExec=run.sh", "test.container")
	diags := qsr002(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	diag := diags[0]
	if diag.Message != "Line is unfinished" {
		t.Errorf("Unexpected diagnostic message: %s", diag.Message)
	}
	if diag.Range.Start.Line != 0 {
		t.Errorf("Expected diagnostic on line 0, got line %d", diag.Range.Start.Line)
	}
	if diag.Source == nil || *diag.Source != "quadlet-lsp.qsr002" {
		t.Errorf("Unexpected diagnostic source: %v", diag.Source)
	}
}

func TestQSR002_CompleteLinesOnly(t *testing.T) {
	s := NewSyntaxChecker("Name=web\nExec=run.sh", "test.container")
	diags := qsr002(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}

func TestQSR002_EqualInValue(t *testing.T) {
	s := NewSyntaxChecker("Env=FOO=bar", "test.container")
	diags := qsr002(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}
