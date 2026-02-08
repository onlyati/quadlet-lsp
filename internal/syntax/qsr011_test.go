package syntax

import (
	"os"
	"path"
	"sync"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCommanderQSR011 struct{}

func (m mockCommanderQSR011) Run(name string, args ...string) ([]string, error) {
	if args[2] == "mock1" {
		return []string{
			`[`,
			`	{`,
			`		 "Config": {`,
			`			"ExposedPorts": {`,
			`				"8080/tcp": {}`,
			`			}`,
			`		 }`,
			`	}`,
			`]`,
		}, nil
	}
	if args[2] == "mock2" {
		return []string{
			`[`,
			`	{`,
			`		 "Config": {`,
			`			"ExposedPorts": {`,
			`				"69/tcp": {}`,
			`			}`,
			`		 }`,
			`	}`,
			`]`,
		}, nil
	}

	return []string{}, nil
}

func TestQSR011_ValidContainer(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock1\nPublishPort=42069:8080")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=mock1\nPublishPort=42069:8080",
			"file://"+tmpDir+"/test1.container",
		),
	}

	for _, s := range cases {
		s.commander = mockCommanderQSR011{}
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr011(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR011_InvalidContainer(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock1\nPublishPort=42069:8081")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=mock1\nPublishPort=42069:8081",
			"file://"+tmpDir+"/test1.container",
		),
	}

	for _, s := range cases {
		s.commander = mockCommanderQSR011{}
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr011(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr011", *diags[0].Source)
		assert.Equal(t, "Port is not exposed in the image, exposed ports: [8080]", diags[0].Message)
	}
}

func TestQSR011_ValidPod(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test.pod", "[Pod]\nPublishPort=42069:8080")
	testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nPod=test.pod\nImage=mock1")
	testutils.CreateTempFile(t, tmpDir, "test2.container", "[Container]\nPod=test.pod\nImage=mock2")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Pod]\nPublishPort=42069:8080",
			"file://"+tmpDir+"/test.pod",
		),
	}

	for _, s := range cases {
		s.commander = mockCommanderQSR011{}
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr011(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR011_InvalidPod(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test.pod", "[Pod]\nPublishPort=42069:5432")
	testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nPod=test.pod\nImage=mock1")
	testutils.CreateTempFile(t, tmpDir, "test2.container", "[Container]\nPod=test.pod\nImage=mock2")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Pod]\nPublishPort=42069:5432",
			"file://"+tmpDir+"/test.pod",
		),
	}

	for _, s := range cases {
		s.commander = mockCommanderQSR011{}
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr011(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr011", *diags[0].Source)
		assert.Equal(t, "Port is not exposed in the image, exposed ports: [8080 69]", diags[0].Message)
	}
}

func TestQSR011_InvalidDropins(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nPublishPort=420:69")
	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, tmpDir+"/foo.container.d", "image.conf", "[Container]\nImage=mock1")

	s := NewSyntaxChecker(
		"[Container]\nPublishPort=420:69",
		"file://"+tmpDir+"/foo.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Mu:            sync.RWMutex{},
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	s.commander = mockCommanderQSR011{}

	diags := qsr011(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr011", *diags[0].Source)
	assert.Equal(t, "Port is not exposed in the image, exposed ports: [8080]", diags[0].Message)
}

func TestQSR011_InvalidMultiDropins(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	testutils.CreateTempDir(t, tmpDir, "foo-bar-baz.container.d")
	testutils.CreateTempDir(t, tmpDir, "foo-bar-.container.d")
	testutils.CreateTempFile(t, tmpDir, "foo-bar-baz.container", "[Container]")
	testutils.CreateTempFile(t, path.Join(tmpDir, "/foo-bar-baz.container.d"), "port.conf", "[Container]\nPublishPort=8080:8080")
	testutils.CreateTempFile(t, path.Join(tmpDir, "/foo-bar-.container.d"), "image.conf", "[Container]\nImage=mock2")

	s := NewSyntaxChecker(
		"[Container]",
		"file://"+tmpDir+"/foo-bar-baz.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
		Mu: sync.RWMutex{},
	}
	s.commander = mockCommanderQSR011{}

	diags := qsr011(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr011", *diags[0].Source)
	assert.Equal(t, "Port is not exposed in the image, exposed ports: [69]", diags[0].Message)
}

func TestQSR011_ValidPodDropins(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo-bar-.container.d")
	testutils.CreateTempFile(t, tmpDir, "foo.pod", "[Pod]\nPublishPort=10080:8080")
	testutils.CreateTempFile(t, tmpDir, "foo-bar-baz.container", "[Container]\nPod=foo.pod")
	testutils.CreateTempFile(t, path.Join(tmpDir, "/foo-bar-.container.d"), "image.conf", "[Container]\nImage=mock1")

	s := NewSyntaxChecker(
		"[Pod]\nPublishPort=10080:8080",
		"file://"+tmpDir+"/foo.pod",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
		Mu: sync.RWMutex{},
	}
	s.commander = mockCommanderQSR011{}

	diags := qsr011(s)
	require.Len(t, diags, 0)
}

func TestQSR011_MoreOption(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock1\nPublishPort=42069:8080/tcp")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=mock1\nPublishPort=42069:8080/tcp",
			"file://"+tmpDir+"/test1.container",
		),
	}

	for _, s := range cases {
		s.commander = mockCommanderQSR011{}
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr011(s)
		require.Len(t, diags, 0)
	}
}
