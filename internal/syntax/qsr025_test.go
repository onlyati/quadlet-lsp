package syntax

import (
	"os"
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.Len(t, diags, 0)
}

func TestQSR025_ValidDropins(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nLabel=app=foo")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.container.d"), "image.conf", "[Container]\nImage=foo.image")

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
	require.Len(t, diags, 0)
}

func TestQSR025_ValidNestedDropins(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo-app")
	testutils.CreateTempDir(t, path.Join(tmpDir, "foo-app"), "foo.container.d")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo-app"), "foo.container", "[Container]\nLabel=app=foo")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo-app", "foo.container.d"), "image.conf", "[Container]\nImage=foo.image")

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
	require.Len(t, diags, 0)
}

func TestQSR025_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nLabel=app=foo")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.container.d"), "volume.conf", "[Container]\nVolume=foo.volume:/app")

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
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr025", *diags[0].Source)
	assert.Equal(t, "Container Quadlet file does not have Image property", diags[0].Message)
}
