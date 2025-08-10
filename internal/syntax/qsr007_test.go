package syntax

import (
	"strings"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR007_Valid(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO1=BAR\nEnvironment=FOO2=BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Build]\nEnvironment=FOO1=BAR\nEnvironment=FOO2=BAR",
			"test.build",
		),
		NewSyntaxChecker(
			`[Container]\nEnvironment=FOO=BAR "MyVar=MyValue" 'foo=bar'\n`,
			"test2.build",
		),
		NewSyntaxChecker(
			`[Container]\nEnvironment=FOO=\n`,
			"test3.build",
		),
		NewSyntaxChecker(
			`[Container]\nEnvironment='fooVariable=barValue'\n`,
			"test4.build",
		),
	}

	for _, s := range variants {
		s.config = &utils.QuadletConfig{}
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
	}

	for _, s := range variants {
		s.config = &utils.QuadletConfig{}
		diags := qsr007(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnosis, but got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr007" {
			t.Fatalf("Exptected quadlet-lsp.qsr007 source but got %s", *diags[0].Source)
		}

		if diags[0].Message != "Invalid format: bad delimiter usage at FOO" {
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
		s.config = &utils.QuadletConfig{}
		diags := qsr007(s)

		if len(diags) == 0 {
			t.Fatalf("Exptected more diagnosis, but got %d at %s", len(diags), s.uri)
		}

		messageCheck := strings.HasPrefix(diags[0].Message, "Invalid format: bad delimiter usage at")
		if !messageCheck {
			t.Fatalf("Got unexpected error message: '%s' at %s", diags[0].Message, s.uri)
		}

	}
}

func TestQSR007_ValidAfter560(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO",
			"test1.container",
		),
	}

	for _, s := range variants {
		s.config = &utils.QuadletConfig{}
		s.config.Podman = utils.BuildPodmanVersion(5, 6, 0)
		diags := qsr007(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnosis, but got %d at %s", len(diags), s.uri)
		}
	}
}
