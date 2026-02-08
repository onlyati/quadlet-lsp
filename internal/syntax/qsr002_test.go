package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR002_UnfinishedLine(t *testing.T) {
	s := NewSyntaxChecker("Name=\nExec=run.sh", "test.container")
	diags := qsr002(s)

	require.Len(t, diags, 1)
	assert.Equal(t, "Line is unfinished", diags[0].Message)
	assert.NotEqual(t, 0, diags[0].Range.Start.Line)

	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr002", *diags[0].Source)
}

func TestQSR002_CompleteLinesOnly(t *testing.T) {
	s := NewSyntaxChecker("Name=web\nExec=run.sh", "test.container")
	diags := qsr002(s)
	require.Len(t, diags, 0)
}

func TestQSR002_EqualInValue(t *testing.T) {
	s := NewSyntaxChecker("Env=FOO=bar", "test.container")
	diags := qsr002(s)
	require.Len(t, diags, 0)
}
