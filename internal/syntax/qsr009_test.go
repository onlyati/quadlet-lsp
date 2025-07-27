package syntax

import "testing"

func TestQSR009_Valid(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nLabel=FOO1=BAR\nLabel=FOO2=BAR",
			"test.container",
		),
		NewSyntaxChecker(
			"[Build]\nLabel=FOO1=BAR\nLabel=FOO2=BAR",
			"test.build",
		),
	}

	for _, s := range variants {
		diags := qsr009(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnosis, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR009_InvalidUnfinished(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nLabel=FOO",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=FOO=",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr009(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnosis, but got %d at %s", len(diags), s.uri)
		}

		if diags[0].Message != "Invalid format of Label specification" {
			t.Fatalf("Got unexpected error message: '%s' at %s", diags[0].Message, s.uri)
		}

	}
}

func TestQSR009_InvalidSpaceFound(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nLabel=FOO = BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=FOO =",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr009(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnosis, but got %d at %s", len(diags), s.uri)
		}

		if diags[0].Message != "Invalid format of Label specification" {
			t.Fatalf("Got unexpected error message: '%s' at %s", diags[0].Message, s.uri)
		}

	}
}
