package syntax

import (
	"errors"
	"path"
	"sync"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
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
	if args[2] == "mock3" {
		return []string{
			`[`,
			`	{`,
			`		 "Config": {`,
			`			"ExposedPorts": {`,
			`				"8080/tcp": {},`,
			`				"8081/tcp": {},`,
			`				"8082/tcp": {},`,
			`				"8083/tcp": {}`,
			`			}`,
			`		 }`,
			`	}`,
			`]`,
		}, nil
	}

	return []string{}, errors.New("invalid image")
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

func TestQSR011_MissingImage(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock0\nPublishPort=42069:8080")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=mock0\nPublishPort=42069:8080",
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
		assert.Equal(t, "Not able to verify exposed ports, because image not pulled: [mock0]", diags[0].Message)
		assert.Equal(t, protocol.DiagnosticSeverityInformation, *diags[0].Severity)
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
		assert.Equal(t, "Port (8081) is not exposed in the image, exposed ports: [8080]", diags[0].Message)
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
		assert.Equal(t, "Port (5432) is not exposed in the image, exposed ports: [8080 69]", diags[0].Message)
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
	assert.Equal(t, "Port (69) is not exposed in the image, exposed ports: [8080]", diags[0].Message)
}

func TestQSR011_InvalidMultiDropins(t *testing.T) {
	tmpDir := t.TempDir()

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
	assert.Equal(t, "Port (8080) is not exposed in the image, exposed ports: [69]", diags[0].Message)
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

func TestQSR011_ValidPortRange(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock3\nPublishPort=42069-42072:8080-8084")

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=mock3\nPublishPort=42080-42083:8080-8083",
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

func TestQSR011_InvalidPortRange(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Invalid start range", func(t *testing.T) {
		testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock3\nPublishPort=42069-42072:hello-8083")

		cases := []SyntaxChecker{
			NewSyntaxChecker(
				"[Container]\nImage=mock3\nPublishPort=42080-42083:hello-8083",
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
			assert.Equal(t, "Not able to verify exposed ports, because start port is not a number: hello", diags[0].Message)
		}
	})
	t.Run("Invalid end range", func(t *testing.T) {
		testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock3\nPublishPort=42069-42072:8080-end")

		cases := []SyntaxChecker{
			NewSyntaxChecker(
				"[Container]\nImage=mock3\nPublishPort=42080-42083:8080-world",
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
			assert.Equal(t, "Not able to verify exposed ports, because end port is not a number: world", diags[0].Message)
		}
	})
	t.Run("Invalid range", func(t *testing.T) {
		testutils.CreateTempFile(t, tmpDir, "test1.container", "[Container]\nImage=mock3\nPublishPort=42069-42072:8080-")

		cases := []SyntaxChecker{
			NewSyntaxChecker(
				"[Container]\nImage=mock3\nPublishPort=42080-42083:8080-",
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
			assert.Equal(t, "Not able to verify exposed ports, because end port is not a number: ", diags[0].Message)
		}
	})
}
