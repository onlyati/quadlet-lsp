package syntax

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR006_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Container]\nImage=foo.image", "foo.container")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	diags := qsr006(s)
	require.Len(t, diags, 0)
}

func TestQSR006_ValidVolume(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo"), "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Volume]\nImage=foo.image", "foo.volume")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	diags := qsr006(s)
	require.Len(t, diags, 0)
}

func TestQSR006_ValidVolumeNested(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Volume]\nImage=foo.image", "foo.volume")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	diags := qsr006(s)
	require.Len(t, diags, 0)
}

func TestQSR006_Skipped(t *testing.T) {
	tmpDir := t.TempDir()

	s := NewSyntaxChecker("[Container]\nImage=library/debian", "foo.container")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	diags := qsr006(s)
	require.Len(t, diags, 0)
}

func TestQSR006_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.image", "[Image]\nImage=docker.io/library/debian")

	s := NewSyntaxChecker("[Container]\nImage=bar.image", "foo.container")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	diags := qsr006(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr006", *diags[0].Source)
	assert.Equal(t, "Image file does not exists: bar.image", diags[0].Message)
}
