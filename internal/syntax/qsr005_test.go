package syntax

import (
	"testing"
)

func TestQSR005_Valid(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nImage=cr.io/org/cont\nAutoUpdate=registry",
		"test.container",
	)

	diags := qsr005(s)

	if len(diags) != 0 {
		t.Fatalf("Exptected 0, but got %d", len(diags))
	}
}

func TestQSR005_ValidKube(t *testing.T) {
	s := NewSyntaxChecker(
		"[Kube]\nYAML=test.yaml\nAutoUpdate=local",
		"test.kube",
	)

	diags := qsr005(s)

	if len(diags) != 0 {
		t.Fatalf("Exptected 0, but got %d", len(diags))
	}
}

func TestQSR005_Invalid(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nImage=cr.io/org/cont\nAutoUpdate=foo",
		"test.container",
	)

	diags := qsr005(s)

	if len(diags) != 1 {
		t.Fatalf("Exptected 1, but got %d", len(diags))
	}

	if diags[0].Message != "Invalid value of AutoUpdate: foo" {
		t.Fatalf("Returned with wrong message tahn expected: %s", diags[0].Message)
	}

	if *diags[0].Source != "quadlet-lsp.qsr005" {
		t.Fatalf("Exptected 'quadlet-lsp.qsr005' source, got %s", *diags[0].Source)
	}
}
