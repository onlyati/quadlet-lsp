package syntax

import "testing"

func TestQSR019_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr019(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR019_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod\nNetwork=my.network",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr019(s)

		if len(diags) != 1 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr019" {
			t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Container cannot have Network because belongs to a pod: test.pod" {
			t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}
