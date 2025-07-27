package syntax

import "testing"

func TestQSR018_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr018(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR018_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod\nPublishPort=420:69",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr018(s)

		if len(diags) != 1 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr018" {
			t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Container cannot have PublishPort because belongs to a pod: test.pod" {
			t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}
