package syntax

import (
	"os"
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR025_ValidBase(t *testing.T) {
	tmpDir := t.TempDir()
	s := NewSyntaxChecker(
		"[Container]\nImage=foo.image",
		"foo.container")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	diags := qsr025(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}
}

func TestQSR025_ValidDropins(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

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
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	diags := qsr025(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}
}

func TestQSR025_ValidNestedDropins(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempDir(
		t,
		tmpDir,
		"foo-app",
	)

	createTempFile(
		t,
		path.Join(tmpDir, "foo-app"),
		"foo.container",
		"[Container]\nLabel=app=foo",
	)
	createTempDir(
		t,
		path.Join(tmpDir, "foo-app"),
		"foo.container.d",
	)
	createTempFile(
		t,
		path.Join(tmpDir, "foo-app", "foo.container.d"),
		"image.conf",
		"[Container]\nImage=foo.image",
	)

	s := NewSyntaxChecker(
		"[Container]\nLabel=app=foo",
		"file://"+path.Join(tmpDir, "foo-app", "foo-container"))
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	diags := qsr025(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, but got %d", len(diags))
	}
}

func TestQSR025_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

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
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

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
