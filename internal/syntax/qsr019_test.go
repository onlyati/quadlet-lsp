package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR019_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod",
			"test1.container",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr019(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR019_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod\nNetwork=my.network",
			"test1.container",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr019(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr019", *diags[0].Source)
		assert.Equal(t, "Container cannot have Network because belongs to a pod: test.pod", diags[0].Message)
	}
}
