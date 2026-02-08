package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR024_Valid(t *testing.T) {
	s := SyntaxChecker{
		documentText: "[Container]\nImage=foo.image\n\n[Service]\nRestart=on-failure\nUser=foo\nGroup=bar\nDynamicUser=true\n",
		uri:          "test.container",
	}

	diags := qsr024(s)
	require.Len(t, diags, 3)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr024", *diags[0].Source)
	assert.Equal(t, "Usage in rootless podman is not recommended: Service.User", diags[0].Message)

	require.NotNil(t, diags[1].Source)
	assert.Equal(t, "quadlet-lsp.qsr024", *diags[1].Source)
	assert.Equal(t, "Usage in rootless podman is not recommended: Service.Group", diags[1].Message)

	require.NotNil(t, diags[2].Source)
	assert.Equal(t, "quadlet-lsp.qsr024", *diags[2].Source)
	assert.Equal(t, "Usage in rootless podman is not recommended: Service.DynamicUser", diags[2].Message)
}
