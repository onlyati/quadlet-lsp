package syntax

import "testing"

func TestQSR007_Valid(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO1=BAR\nEnvironment=FOO2=BAR",
			"test.container",
		),
		NewSyntaxChecker(
			"[Build]\nEnvironment=FOO1=BAR\nEnvironment=FOO2=BAR",
			"test.build",
		),
	}

	for _, s := range variants {
		diags := qsr007(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnosis, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR007_InvalidUnfinished(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO=",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr007(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnosis, but got %d at %s", len(diags), s.uri)
		}

		if diags[0].Message != "Invalid format of Environment variable specification" {
			t.Fatalf("Got unexpected error message: '%s' at %s", diags[0].Message, s.uri)
		}

	}
}

func TestQSR007_InvalidSpaceFound(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO = BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO =",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr007(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnosis, but got %d at %s", len(diags), s.uri)
		}

		if diags[0].Message != "Invalid format of Environment variable specification" {
			t.Fatalf("Got unexpected error message: '%s' at %s", diags[0].Message, s.uri)
		}

	}
}
