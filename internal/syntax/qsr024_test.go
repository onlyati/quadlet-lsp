package syntax

import (
	"strings"
	"testing"
)

func TestQSR024_Valid(t *testing.T) {
	s := SyntaxChecker{
		documentText: "[Container]\nImage=foo.image\n\n[Service]\nRestart=on-failure\nUser=foo\nGroup=bar\nDynamicUser=true\n",
		uri:          "test.container",
	}

	diags := qsr024(s)

	if len(diags) != 3 {
		t.Fatalf("expected 3 diagnostics but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr024" {
		t.Fatalf("unexpected source: %s", *diags[0].Source)
	}

	if !strings.HasPrefix(diags[0].Message, "Usage in rootless podman is not recommended:") {
		t.Fatalf("unexpected message: %s", diags[0].Message)
	}
}
