package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR013_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "data1.volume", "[Volume]")
	testutils.CreateTempFile(t, tmpDir, "data2.volume", "[Volume]")
	testutils.CreateTempFile(t, tmpDir, "data@.volume", "[Volume]")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.build",
		),
		NewSyntaxChecker(
			"[Build]\nVolume=data@%i.volume:/app:r\nVolume=data@test.volume:/data/:rw",
			"file://"+tmpDir+"/test2.build",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr013(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR013_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.build",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr013(s)
		require.Len(t, diags, 2)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr013", *diags[0].Source)
		assert.Equal(t, "Volume file does not exists: data1.volume", diags[0].Message)

		require.NotNil(t, diags[1].Source)
		assert.Equal(t, "quadlet-lsp.qsr013", *diags[1].Source)
		assert.Equal(t, "Volume file does not exists: data2.volume", diags[1].Message)

	}
}
