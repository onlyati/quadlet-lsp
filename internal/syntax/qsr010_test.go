package syntax

import "testing"

func TestQSR010_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:420:69",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nPublishPort=10.0.0.1:420:69",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=10.0.0.1:420:69",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=420:69",
			"test4.container",
		),
		NewSyntaxChecker(
			"[Pod]\nPublishPort=420:69",
			"test5.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=420:69",
			"test6.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=:69",
			"test7.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=10.0.0.1::69",
			"test8.container",
		),
	}

	for _, s := range cases {
		d := qsr010(s)

		if len(d) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(d), s.uri)
		}
	}
}

func TestQSR010_InvalidFormat(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=420",
			"test2.container",
		),
	}

	for _, s := range cases {
		d := qsr010(s)

		if len(d) != 1 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(d), s.uri)
		}

		if *d[0].Source != "quadlet-lsp.qsr010" {
			t.Fatalf("Wrong source, expected 'quadlet-lsp.qsr010', got '%s' at %s", *d[0].Source, s.uri)
		}

		if d[0].Message != "Incorrect format of PublishPort: invalid format" {
			t.Fatalf("Unexpected message: '%s' at %s", d[0].Message, s.uri)
		}
	}
}

func TestQSR010_InvalidPortIsText(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:nice:420",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:69:ez",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=nice:420",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=69:ez",
			"test4.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=69:",
			"test4.container",
		),
	}

	for _, s := range cases {
		d := qsr010(s)

		if len(d) != 1 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(d), s.uri)
		}

		if *d[0].Source != "quadlet-lsp.qsr010" {
			t.Fatalf("Wrong source, expected 'quadlet-lsp.qsr010', got '%s' at %s", *d[0].Source, s.uri)
		}

		if d[0].Message != "Incorrect format of PublishPort: not a number" {
			t.Fatalf("Unexpected message: '%s' at %s", d[0].Message, s.uri)
		}
	}
}

func TestQSR010_InvalidInvalidPortNumber(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:-69:420",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:69:80000",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=-69:420",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=69:80000",
			"test4.container",
		),
	}

	for _, s := range cases {
		d := qsr010(s)

		if len(d) != 1 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(d), s.uri)
		}

		if *d[0].Source != "quadlet-lsp.qsr010" {
			t.Fatalf("Wrong source, expected 'quadlet-lsp.qsr010', got '%s' at %s", *d[0].Source, s.uri)
		}

		if d[0].Message != "Incorrect format of PublishPort: port must be between [0;65535]" {
			t.Fatalf("Unexpected message: '%s' at %s", d[0].Message, s.uri)
		}
	}
}
