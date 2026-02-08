package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR014_Valid(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "net1.network", "[Network]")
	testutils.CreateTempFile(t, tmpDir, "net2.network", "[Network]")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.build",
		),
		NewSyntaxChecker(
			"[Kube]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.kube",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr014(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR014_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.build",
		),
		NewSyntaxChecker(
			"[Kube]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.kube",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr014(s)
		require.Len(t, diags, 2)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr014", *diags[0].Source)
		assert.Equal(t, "Network file does not exists: net1.network", diags[0].Message)

		require.NotNil(t, diags[1].Source)
		assert.Equal(t, "quadlet-lsp.qsr014", *diags[1].Source)
		assert.Equal(t, "Network file does not exists: net2.network", diags[1].Message)
	}
}
