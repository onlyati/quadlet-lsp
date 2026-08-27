package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR027(t *testing.T) {
	t.Run("Valid 1", func(t *testing.T) {
		s := SyntaxChecker{
			documentText: "[Container]\nImageVolume=bind\n",
			uri:          "test.container",
		}
		s.config = &utils.QuadletConfig{
			Podman: utils.BuildPodmanVersion(6, 1, 0),
		}
		diags := qsr027(s)
		require.Len(t, diags, 0)
	})
	t.Run("Valid 2", func(t *testing.T) {
		s := SyntaxChecker{
			documentText: "[Container]\nImageVolume=tmpfs\n",
			uri:          "test.container",
		}
		s.config = &utils.QuadletConfig{
			Podman: utils.BuildPodmanVersion(6, 1, 0),
		}
		diags := qsr027(s)
		require.Len(t, diags, 0)
	})
	t.Run("Valid 3", func(t *testing.T) {
		s := SyntaxChecker{
			documentText: "[Container]\nImageVolume=ignore\n",
			uri:          "test.container",
		}
		s.config = &utils.QuadletConfig{
			Podman: utils.BuildPodmanVersion(6, 1, 0),
		}
		diags := qsr027(s)
		require.Len(t, diags, 0)
	})
	t.Run("InValid value", func(t *testing.T) {
		s := SyntaxChecker{
			documentText: "[Container]\nImageVolume=foo\n",
			uri:          "test.container",
		}
		s.config = &utils.QuadletConfig{
			Podman: utils.BuildPodmanVersion(6, 1, 0),
		}
		diags := qsr027(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr027", *diags[0].Source)
		assert.Equal(t, "Allowed values are bind, tmpfs and ignore, it is: foo", diags[0].Message)
	})
}
