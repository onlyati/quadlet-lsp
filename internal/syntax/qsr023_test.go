package syntax

import (
	"strings"
	"testing"
)

func TestQSR023_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Unit]\nWants=%N-db.container",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=Unit=%n",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=%h/app1:/app:ro",
			"test3.container",
		),
	}

	for _, s := range cases {
		diags := qsr023(s)

		if len(diags) != 0 {
			t.Fatalf("expected 0 finding, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR023_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Unit]\nWants=%r-db.container",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=Unit=%ö",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=%5/app1:/app:ro",
			"test3.pod",
		),
	}

	for _, s := range cases {
		diags := qsr023(s)

		if len(diags) != 1 {
			t.Fatalf("expected 1 finding, but got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr023" {
			t.Fatalf("unexpected source: %s at %s", *diags[0].Source, s.uri)
		}

		msgStart := strings.HasPrefix(diags[0].Message, "Specifier ")
		msgEnd := strings.HasSuffix(diags[0].Message, "is invalid")
		if !msgStart || !msgEnd {
			t.Fatalf("unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}

		if diags[0].Range.Start.Line != 1 {
			t.Fatalf("found issue in unexpteced line: %d at %s", diags[0].Range.Start.Line, s.uri)
		}
	}
}
