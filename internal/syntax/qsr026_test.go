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

func TestQSR026_ValidBase(t *testing.T) {
	s := NewSyntaxChecker(
		"[Artifact]\nArtifact=foo.io/bar/example:latest",
		"foo.container")
	s.config = &utils.QuadletConfig{}

	diags := qsr026(s)
	require.Len(t, diags, 0)
}

func TestQSR026_ValidDropins(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo.artifact.d")
	testutils.CreateTempFile(t, tmpDir, "foo.artifact", "[Artifact]\nAuthFile=/etc/registry/auth.json")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.artifact.d"), "artifact.conf", "[Artifact]\nArtifact=foo.io/bar/example2:latest")

	s := NewSyntaxChecker(
		"[Artifact]\nAuthFile=/etc/registry/auth.json",
		"file://"+tmpDir+string(os.PathSeparator)+"foo.artifact")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	diags := qsr026(s)
	require.Len(t, diags, 0)
}

func TestQSR026_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo.artifact.d")
	testutils.CreateTempFile(t, tmpDir, "foo.artifact", "[Artifact]\nAuthFile=/etc/registry/auth.json")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.artifact.d"), "volume.conf", "[Artifact]\nTLSVerify=false")

	s := NewSyntaxChecker(
		"[Artifact]\nAuthFile=/etc/registry/auth.json",
		"file://"+tmpDir+string(os.PathSeparator)+"foo.artifact")
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	diags := qsr026(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr026", *diags[0].Source)
	assert.Equal(t, "Artifact Quadlet file does not have Artifact property", diags[0].Message)
}
