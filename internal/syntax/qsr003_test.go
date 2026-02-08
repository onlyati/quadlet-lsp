package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR003_ValidProperties(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nContainerName=app\nExec=run.sh\n# Something=comment\nUser=root",
		"test.container",
	)
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 5, 2),
	}
	diags := qsr003(s)
	require.Len(t, diags, 0)
}

func TestQSR003_InvalidProperty(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nContainerName=app\nFoobar=yes\nExec=run.sh", "test.container")
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 5, 2),
	}
	diags := qsr003(s)

	require.Len(t, diags, 1)
	assert.Equal(t, "Invalid property is found: Container.Foobar", diags[0].Message)
	assert.Equal(t, uint32(2), diags[0].Range.Start.Line)

	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr003", *diags[0].Source)
}

func TestQSR003_UnknownSection(t *testing.T) {
	s := NewSyntaxChecker("[Test]\nDescription=42", "test.container")
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 5, 2),
	}
	diags := qsr003(s)
	require.Len(t, diags, 2)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr003", *diags[0].Source)
	assert.Equal(t, "Invalid property is found: Test.", diags[0].Message)

	require.NotNil(t, diags[1].Source)
	assert.Equal(t, "Invalid property is found: Test.Description", diags[1].Message)
	assert.Equal(t, "quadlet-lsp.qsr003", *diags[1].Source)
}

func TestQSR003_OldVersion(t *testing.T) {
	// Memory for container is available from 5.5.0
	s := NewSyntaxChecker("[Container]\nContainerName=app\nMemory=512M", "test.container")
	s.config = &utils.QuadletConfig{
		Podman: utils.BuildPodmanVersion(5, 4, 2),
	}
	diags := qsr003(s)

	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr003", *diags[0].Source)
	assert.Equal(t, "Invalid property is found: Container.Memory", diags[0].Message)
}
