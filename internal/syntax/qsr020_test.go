package syntax

import (
	"strings"
	"testing"
)

func TestQSR020_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nContainerName=foo\n",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Volume]\nVolumeName=foo\n",
			"test1.volume",
		),
		NewSyntaxChecker(
			"[Network]\nNetworkName=foo\n",
			"test1.network",
		),
		NewSyntaxChecker(
			"[Pod]\nPodName=foo\n",
			"test1.pod",
		),
	}

	for _, s := range cases {
		diags := qsr020(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnostics but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR020_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nContainerName=.foo\n",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Volume]\nVolumeName=_foo\n",
			"test1.volume",
		),
		NewSyntaxChecker(
			"[Network]\nNetworkName=-foo\n",
			"test1.network",
		),
		NewSyntaxChecker(
			"[Pod]\nPodName=*foo\n",
			"test1.pod",
		),
	}

	for _, s := range cases {
		diags := qsr020(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnostics but got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr020" {
			t.Fatalf("Unexpected source: %s", *diags[0].Source)
		}

		matchMessageStart := strings.HasPrefix(diags[0].Message, "Invalid name of unit: ")
		if !matchMessageStart {
			t.Fatalf("Unexpected error message: %s", diags[0].Message)
		}
	}
}
