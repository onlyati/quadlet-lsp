package syntax

import (
	"os"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR025_ValidBase(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nImage=foo.image",
		"foo.container")
	s.config = &utils.QuadletConfig{}

	diags := qsr025(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}
}

func TestQSR025_ValidDropins(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t, tmpDir,
		"foo.container",
		"[Container]\nLabel=app=foo",
	)
	createTempDir(
		t, tmpDir,
		"foo.container.d",
	)
	createTempFile(
		t,
		tmpDir+string(os.PathSeparator)+"foo.container.d",
		"image.conf",
		"[Container]\nImage=foo.image",
	)

	s := NewSyntaxChecker(
		"[Container]\nLabel=app=foo",
		"file://"+tmpDir+string(os.PathSeparator)+"foo.container")
	s.config = &utils.QuadletConfig{}

	diags := qsr025(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}
}

func TestQSR025_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t, tmpDir,
		"foo.container",
		"[Container]\nLabel=app=foo",
	)
	createTempDir(
		t, tmpDir,
		"foo.container.d",
	)
	createTempFile(
		t,
		tmpDir+string(os.PathSeparator)+"foo.container.d",
		"volume.conf",
		"[Container]\nVolume=foo.volume:/app",
	)

	s := NewSyntaxChecker(
		"[Container]\nLabel=app=foo",
		"file://"+tmpDir+string(os.PathSeparator)+"foo.container")
	s.config = &utils.QuadletConfig{}

	diags := qsr025(s)

	if len(diags) != 1 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}

	if *diags[0].Source != "quadlet-lsp.qsr025" {
		t.Fatalf("expected quadlet-lsp.qsr025 source, but got '%s'", *diags[0].Source)
	}

	if diags[0].Message != "Container Quadlet file does not have Image property" {
		t.Fatalf("Unpextected message: '%s'", diags[0].Message)
	}
}
