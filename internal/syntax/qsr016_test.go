package syntax

import "testing"

func TestQSR016_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nUserNS=keep-id",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=keep-id:uid=101,gid=101",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=auto",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=host",
			"test4.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=nomap",
			"test5.container",
		),
	}

	for _, s := range cases {
		diags := qsr015(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR016_NoParameters(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nUserNS=auto:gid=101",
		"test.container",
	)

	diags := qsr016(s)

	if len(diags) != 1 {
		t.Fatalf("Exptected 0 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr016" {
		t.Fatalf("Exptected quadlet-lsp.qsr016 source, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid value of UserNS: 'auto' has no parameters" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}

func TestQSR016_KeepIdWrongParameter(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nUserNS=keep-id:gid=101,foo=101",
		"test.container",
	)

	diags := qsr016(s)

	if len(diags) != 1 {
		t.Fatalf("Exptected 0 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr016" {
		t.Fatalf("Exptected quadlet-lsp.qsr016 source, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid value of UserNS: [uid gid] allowed but found foo=101" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}

func TestQSR016_InvalidValue(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nUserNS=foo",
		"test.container",
	)

	diags := qsr016(s)

	if len(diags) != 1 {
		t.Fatalf("Exptected 0 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr016" {
		t.Fatalf("Exptected quadlet-lsp.qsr016 source, but got %s", *diags[0].Source)
	}

	if diags[0].Message != "Invalid value of UserNS: allowed values: '[auto host keep-id nomap]' and found foo" {
		t.Fatalf("Unexpected message: %s", diags[0].Message)
	}
}
