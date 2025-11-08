package syntax

import (
	"os"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR026_ValidBase(t *testing.T) {
	s := NewSyntaxChecker(
		"[Artifact]\nArtifact=foo.io/bar/example:latest",
		"foo.container")
	s.config = &utils.QuadletConfig{}

	diags := qsr026(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}
}

func TestQSR026_ValidDropins(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t, tmpDir,
		"foo.artifact",
		"[Artifact]\nAuthFile=/etc/registry/auth.json",
	)
	createTempDir(
		t, tmpDir,
		"foo.artifact.d",
	)
	createTempFile(
		t,
		tmpDir+string(os.PathSeparator)+"foo.artifact.d",
		"artifact.conf",
		"[Artifact]\nArtifact=foo.io/bar/example2:latest",
	)

	s := NewSyntaxChecker(
		"[Artifact]\nAuthFile=/etc/registry/auth.json",
		"file://"+tmpDir+string(os.PathSeparator)+"foo.artifact")
	s.config = &utils.QuadletConfig{}

	diags := qsr026(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}
}

func TestQSR026_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t, tmpDir,
		"foo.artifact",
		"[Artifact]\nAuthFile=/etc/registry/auth.json",
	)
	createTempDir(
		t, tmpDir,
		"foo.artifact.d",
	)
	createTempFile(
		t,
		tmpDir+string(os.PathSeparator)+"foo.artifact.d",
		"volume.conf",
		"[Artifact]\nTLSVerify=false",
	)

	s := NewSyntaxChecker(
		"[Artifact]\nAuthFile=/etc/registry/auth.json",
		"file://"+tmpDir+string(os.PathSeparator)+"foo.artifact")
	s.config = &utils.QuadletConfig{}

	diags := qsr026(s)

	if len(diags) != 1 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr026" {
		t.Fatalf("expected quadlet-lsp.qsr026 source, but got '%s'", *diags[0].Source)
	}

	if diags[0].Message != "Artifact Quadlet file does not have Artifact property" {
		t.Fatalf("Unpextected message: '%s'", diags[0].Message)
	}
}
