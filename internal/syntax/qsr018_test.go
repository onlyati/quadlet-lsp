package syntax

import (
	"os"
	"sync"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR018_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"test1.container",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
		}
		diags := qsr018(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR018_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod\nPublishPort=420:69",
			"test1.container",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
		}
		diags := qsr018(s)

		if len(diags) != 1 {
			t.Fatalf("Expected 1 diagnostics, got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr018" {
			t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Container cannot have PublishPort because belongs to a pod: test.pod" {
			t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}

func TestQSR018_InvalidWithDropins(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.container", "[Container]\nImage=foo.image\nPublishPort=8080:8080")
	createTempDir(t, tmpDir, "foo.container.d")
	createTempFile(t, tmpDir+"/foo.container.d", "10-pod.conf", "[Container]\nPod=foo.pod")

	s := NewSyntaxChecker(
		"[Container]\nImage=foo.image\nPublishPort=8080:8080",
		"file:///"+tmpDir+"/foo.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Mu:            sync.RWMutex{},
	}

	diags := qsr018(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostics, got %d at %s", len(diags), s.uri)
	}

	if *diags[0].Source != "quadlet-lsp.qsr018" {
		t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
	}

	if diags[0].Message != "Container cannot have PublishPort because belongs to a pod: foo.pod" {
		t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
	}
}

func TestQSR018_InvalidWithDropinsMoreLevel(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo-bar-baz.container", "[Container]\nImage=foo.image\nPublishPort=8080:8080")
	createTempDir(t, tmpDir, "foo-bar-baz.container.d")
	createTempFile(t, tmpDir+"/foo-bar-baz.container.d", "10-network.conf", "[Container]\nNetwork=foo.network")
	createTempDir(t, tmpDir, "foo-bar-.container.d")
	createTempFile(t, tmpDir+"/foo-bar-.container.d", "10-pod.conf", "[Container]\nPod=foo.pod")

	s := NewSyntaxChecker(
		"[Container]\nImage=foo.image\nPublishPort=8080:8080",
		"file:///"+tmpDir+"/foo-bar-baz.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Mu:            sync.RWMutex{},
	}

	diags := qsr018(s)

	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostics, got %d at %s", len(diags), s.uri)
	}

	if *diags[0].Source != "quadlet-lsp.qsr018" {
		t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
	}

	if diags[0].Message != "Container cannot have PublishPort because belongs to a pod: foo.pod" {
		t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
	}
}
