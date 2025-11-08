package syntax

import "testing"

func TestQSR004_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=docker.io/library/debian:bookworm-slim",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=ghcr.io/henrygd/beszel/beszel-agent:latest",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost/test",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost:5000/test",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=example.com:5000/test",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Artifact]\nArtifact=example.com:5000/test",
			"test3.artifact",
		),
	}

	for _, s := range cases {
		diags := qsr004(s)

		if len(diags) != 0 {
			t.Errorf("Expected no diagnostics, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR004_ValidImageFile(t *testing.T) {
	s := NewSyntaxChecker("[Image]\nImage=docker.io/library/debian:bookworm-slim", "test.image")
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

func TestQSR004_ValidWithBuild(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nImage=db.build", "test.container")
	diags := qsr004(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}

func TestQSR004_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=library/debian:bookworm-slim",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost/",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Artifact]\nArtifact=localhost/",
			"test1.artifact",
		),
	}

	for _, s := range cases {
		diags := qsr004(s)

		if len(diags) != 1 {
			t.Errorf("Expected 1 diagnostic, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR004_InvalidImageFile(t *testing.T) {
	s := NewSyntaxChecker("[Image]\nImage=library/debian:bookworm-slim", "test.image")
	diags := qsr004(s)

	if len(diags) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(diags))
	}
}

func TestQSR004_NonContainer(t *testing.T) {
	s := NewSyntaxChecker("[Pod]\nImage=library/debian:bookworm-slim", "test.pod")
	diags := qsr004(s)

	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics, got %d", len(diags))
	}
}
