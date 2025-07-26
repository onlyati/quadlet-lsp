package syntax

import "testing"

func TestQSR004_Valid(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nImage=docker.io/library/debian:bookworm-slim", "test.container")
	diags := qsr004(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}

func TestQSR004_ValidWithImage(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nImage=db.image", "test.container")
	diags := qsr004(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}

func TestQSR004_Invalid(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nImage=library/debian:bookworm-slim", "test.container")
	diags := qsr004(s)

	if len(diags) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(diags))
	}
}

func TestQSR004_NonContainer(t *testing.T) {
	s := NewSyntaxChecker("[Pod]\nImage=library/debian:bookworm-slim", "test.container")
	diags := qsr004(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}
