package syntax

import (
	"sync"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR018_Valid(t *testing.T) {
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
		diags := qsr018(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR018_Invalid(t *testing.T) {
	tmpDir := t.TempDir()

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPod=test.pod\nPublishPort=420:69",
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
		diags := qsr018(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr018", *diags[0].Source)
		assert.Equal(t, "Container cannot have PublishPort because belongs to a pod: test.pod", diags[0].Message)
	}
}

func TestQSR018_InvalidWithDropins(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nImage=foo.image\nPublishPort=8080:8080")
	testutils.CreateTempFile(t, tmpDir+"/foo.container.d", "10-pod.conf", "[Container]\nPod=foo.pod")

	s := NewSyntaxChecker(
		"[Container]\nImage=foo.image\nPublishPort=8080:8080",
		"file:///"+tmpDir+"/foo.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
		Mu: sync.RWMutex{},
	}

	diags := qsr018(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr018", *diags[0].Source)
	assert.Equal(t, "Container cannot have PublishPort because belongs to a pod: foo.pod", diags[0].Message)
}

func TestQSR018_InvalidWithDropinsMoreLevel(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo-bar-baz.container.d")
	testutils.CreateTempDir(t, tmpDir, "foo-bar-.container.d")
	testutils.CreateTempFile(t, tmpDir, "foo-bar-baz.container", "[Container]\nImage=foo.image\nPublishPort=8080:8080")
	testutils.CreateTempFile(t, tmpDir+"/foo-bar-baz.container.d", "10-network.conf", "[Container]\nNetwork=foo.network")
	testutils.CreateTempFile(t, tmpDir+"/foo-bar-.container.d", "10-pod.conf", "[Container]\nPod=foo.pod")

	s := NewSyntaxChecker(
		"[Container]\nImage=foo.image\nPublishPort=8080:8080",
		"file:///"+tmpDir+"/foo-bar-baz.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
		Mu: sync.RWMutex{},
	}

	diags := qsr018(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr018", *diags[0].Source)
	assert.Equal(t, "Container cannot have PublishPort because belongs to a pod: foo.pod", diags[0].Message)
}
