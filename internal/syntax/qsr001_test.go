package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR001_WithValidSection(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nName=app", "test.container")

	diags := qsr001(s)

	require.Len(t, diags, 0)
}

func TestQSR001_WithoutValidSection(t *testing.T) {
	s := NewSyntaxChecker("Name=app\nExec=run.sh", "test.container")

	diags := qsr001(s)

	require.Len(t, diags, 1)
	assert.Equal(t, "quadlet-lsp.qsr001", *diags[0].Source)

	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "Missing any of these sections: [[Image] [Container] [Volume] [Network] [Kube] [Pod] [Build] [Artifact]]", diags[0].Message)
}
