package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR017_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test.pod", "[Pod]")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"file://"+tmpDir+"/test1.container",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr017(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR017_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"file://"+tmpDir+"/test1.container",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr017(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr017", *diags[0].Source)
		assert.Equal(t, "Pod file does not exists: test.pod", diags[0].Message)
	}
}
