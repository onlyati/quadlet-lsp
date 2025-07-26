package syntax

import "testing"

func TestQSR008_Valid(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO1=BAR\nAnnotation=FOO2=BAR",
			"test.container",
		),
		NewSyntaxChecker(
			"[Build]\nAnnotation=FOO1=BAR\nAnnotation=FOO2=BAR",
			"test.build",
		),
	}

	for _, s := range variants {
		diags := qsr008(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnosis, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR008_InvalidUnfinished(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO=",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr008(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnosis, but got %d at %s", len(diags), s.uri)
		}

		if diags[0].Message != "Invalid format of Annotation specification" {
			t.Fatalf("Got unexpected error message: '%s' at %s", diags[0].Message, s.uri)
		}

	}
}

func TestQSR008_InvalidSpaceFound(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO = BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO =",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr008(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnosis, but got %d at %s", len(diags), s.uri)
		}

		if diags[0].Message != "Invalid format of Annotation specification" {
			t.Fatalf("Got unexpected error message: '%s' at %s", diags[0].Message, s.uri)
		}

	}
}
