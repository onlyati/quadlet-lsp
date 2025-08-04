package syntax

import (
	"strings"
	"testing"
)

func TestQSR012_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]Secret=my-secret",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,target=/app/pw",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,type=mount,target=/app/pw",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,uid=69,gid=420,mode=0400",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,type=mount,uid=69,gid=420,mode=0400",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr012(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR012_UnfinishedOption(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nSecret=my-secret,type\n",
		"foo.container",
	)

	diags := qsr012(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr012" {
		t.Fatalf("Exptected source quadlet-lsp.qsr012, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid format of secret specification: 'type' has no value" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}

func TestQSR012_InvalidType(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nSecret=my-secret,type=foo\n",
		"foo.container",
	)

	diags := qsr012(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr012" {
		t.Fatalf("Exptected source quadlet-lsp.qsr012, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid format of secret specification: 'type' can be either 'mount' or 'env'" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}

func TestQSR012_InvalidOption(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nSecret=my-secret,foo=bar\n",
		"foo.container",
	)

	diags := qsr012(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr012" {
		t.Fatalf("Exptected source quadlet-lsp.qsr012, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid format of secret specification: 'foo' is invalid option" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}

func TestQSR012_InvalidWithEnv(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nSecret=my-secret,type=env,uid=69,gid=420,mode=0400\n",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr012(s)

		if len(diags) != 1 {
			t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
		}

		if *diags[0].Source != "quadlet-lsp.qsr012" {
			t.Fatalf("Exptected source quadlet-lsp.qsr012, but got %s", *diags[0].Source)
		}

		checkMessageStart := strings.HasPrefix(diags[0].Message, "Invalid format of secret specification: ")
		checkMessageSuffix := strings.HasSuffix(diags[0].Message, "' only allowed if type=mount")
		if !checkMessageStart || !checkMessageSuffix {
			t.Fatalf("Unexpected message: %s", diags[0].Message)
		}
	}
}
