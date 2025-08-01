package syntax

import (
	"testing"
)

func TestQSR015_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=foo.volume:/app",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=foo.volume:/app:ro",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=foo.volume:/app:ro,Z",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr015(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR015_InvalidContainerDirectory(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nVolume=foo.volume:data/config",
		"test1.container",
	)

	diags := qsr015(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr015" {
		t.Fatalf("Exptected quadlet-lsp.qsr015 source, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid format of Volume specification: container directory is not absolute" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}

func TestQSR016_UnkownFlag(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nVolume=foo.volume:/app/data/config:rw,Z,foo,nocopy",
		"test1.container",
	)

	diags := qsr015(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr015" {
		t.Fatalf("Exptected quadlet-lsp.qsr015 source, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid format of Volume specification: 'foo' flag is unknown" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}
