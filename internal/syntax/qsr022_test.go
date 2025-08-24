package syntax

import (
	"testing"
)

func TestQSR022_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/container/%n",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:%t",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/%N",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/%",
			"test3.container",
		),
	}

	for _, s := range cases {
		diags := qsr022(s)

		if len(diags) != 0 {
			t.Fatalf("expected 0 finding, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR022_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/container/%t",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/%t",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=/%t:%t",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr022(s)

		if len(diags) != 1 {
			t.Fatalf("expected 1 finding, but got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr022" {
			t.Fatalf("unexpected source: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Specifier, %t, already starts with '/' sign" {
			t.Fatalf("unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}

		if diags[0].Range.Start.Line != 1 {
			t.Fatalf("found issue in unexpteced line: %d at %s", diags[0].Range.Start.Line, s.uri)
		}
	}
}

func TestQSR022_InvalidMultipleInOneLine(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=/%t:/container/%t/test",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr022(s)

		if len(diags) != 2 {
			t.Fatalf("expected 2 finding, but got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr022" {
			t.Fatalf("unexpected source: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Specifier, %t, already starts with '/' sign" {
			t.Fatalf("unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}

		if diags[0].Range.Start.Line != 1 {
			t.Fatalf("found issue in unexpteced line: %d at %s", diags[0].Range.Start.Line, s.uri)
		}

		line := "Volume=/%t:/container/%t"
		problemPart := line[diags[0].Range.Start.Character : diags[0].Range.End.Character+1]
		if problemPart != "/%t" {
			t.Fatalf(
				"expected '/%%t' as problem but got '%s' at %s",
				problemPart,
				s.uri,
			)
		}

		problemPart = line[diags[1].Range.Start.Character : diags[1].Range.End.Character+1]
		if problemPart != "/%t" {
			t.Fatalf(
				"expected '/%%h' as problem but got '%s' at %s",
				problemPart,
				s.uri,
			)
		}
	}
}
