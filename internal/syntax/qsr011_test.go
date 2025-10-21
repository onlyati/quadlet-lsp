package syntax

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

type mockCommanderQSR011 struct{}

func (m mockCommanderQSR011) Run(name string, args ...string) ([]string, error) {
	if args[2] == "mock1" {
		return []string{
			"[",
			"	{",
			"		 \"Config\": {",
			"			\"ExposedPorts\": {",
			"				\"8080/tcp\": {}",
			"			}",
			"		 }",
			"	}",
			"]",
		}, nil
	}
	if args[2] == "mock2" {
		return []string{
			"[",
			"	{",
			"		 \"Config\": {",
			"			\"ExposedPorts\": {",
			"				\"69/tcp\": {}",
			"			}",
			"		 }",
			"	}",
			"]",
		}, nil
	}

	return []string{}, nil
}

func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0o644)
	assert.NoError(t, err)
	return path
}

func createTempDir(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.Mkdir(path, 0o755)
	assert.NoError(t, err)
	return path
}

func TestQSR011_ValidContainer(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"test1.container",
		"[Container]\nImage=mock1\nPublishPort=42069:8080",
	)

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
		}
		diags := qsr011(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR011_InvalidContainer(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"test1.container",
		"[Container]\nImage=mock1\nPublishPort=42069:8081",
	)

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
		}
		diags := qsr011(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnostics, but got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr011" {
			t.Fatalf("Wrong source is used at %s", s.uri)
		}

		if diags[0].Message != "Port is not exposed in the image, exposed ports: [8080]" {
			t.Fatalf("Unexptected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}

func TestQSR011_ValidPod(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"test.pod",
		"[Pod]\nPublishPort=42069:8080",
	)

	createTempFile(
		t,
		tmpDir,
		"test1.container",
		"[Container]\nPod=test.pod\nImage=mock1",
	)

	createTempFile(
		t,
		tmpDir,
		"test2.container",
		"[Container]\nPod=test.pod\nImage=mock2",
	)

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
		}
		diags := qsr011(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR011_InvalidPod(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"test.pod",
		"[Pod]\nPublishPort=42069:5432",
	)

	createTempFile(
		t,
		tmpDir,
		"test1.container",
		"[Container]\nPod=test.pod\nImage=mock1",
	)

	createTempFile(
		t,
		tmpDir,
		"test2.container",
		"[Container]\nPod=test.pod\nImage=mock2",
	)

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
		}
		diags := qsr011(s)

		if len(diags) != 1 {
			t.Fatalf("Exptected 1 diagnostics, but got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr011" {
			t.Fatalf("Wrong source is used at %s", s.uri)
		}

		if diags[0].Message != "Port is not exposed in the image, exposed ports: [8080 69]" {
			t.Fatalf("Unexptected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}

func TestQSR011_InvalidDropins(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.container", "[Container]\nPublishPort=420:69")
	createTempDir(t, tmpDir, "foo.container.d")
	createTempFile(t, tmpDir+"/foo.container.d", "image.conf", "[Container]\nImage=mock1")

	s := NewSyntaxChecker(
		"[Container]\nPublishPort=420:69",
		"file://"+tmpDir+"/foo.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Mu:            sync.RWMutex{},
	}
	s.commander = mockCommanderQSR011{}

	diags := qsr011(s)

	if len(diags) != 1 {
		t.Fatalf("expected 1 diags, but got %d", len(diags))
	}
}

func TestQSR011_InvalidMultiDropins(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo-bar-baz.container", "[Container]")
	createTempDir(t, tmpDir, "foo-bar-baz.container.d")
	createTempFile(t, tmpDir+"/foo-bar-baz.container.d", "port.conf", "[Container]\nPublishPort=8080:8080")
	createTempDir(t, tmpDir, "foo-bar-.container.d")
	createTempFile(t, tmpDir+"/foo-bar-.container.d", "image.conf", "[Container]\nImage=mock2")

	s := NewSyntaxChecker(
		"[Container]",
		"file://"+tmpDir+"/foo-bar-baz.container",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Mu:            sync.RWMutex{},
	}
	s.commander = mockCommanderQSR011{}

	diags := qsr011(s)

	if len(diags) != 1 {
		t.Fatalf("expected 1 diags, but got %d", len(diags))
	}
}

func TestQSR011_ValidPodDropins(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.pod", "[Pod]\nPublishPort=10080:8080")
	createTempFile(t, tmpDir, "foo-bar-baz.container", "[Container]\nPod=foo.pod")
	createTempDir(t, tmpDir, "foo-bar-.container.d")
	createTempFile(t, tmpDir+"/foo-bar-.container.d", "image.conf", "[Container]\nImage=mock1")

	s := NewSyntaxChecker(
		"[Pod]\nPublishPort=10080:8080",
		"file://"+tmpDir+"/foo.pod",
	)
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Mu:            sync.RWMutex{},
	}
	s.commander = mockCommanderQSR011{}

	diags := qsr011(s)

	if len(diags) != 0 {
		t.Fatalf("expected 0 diags, but got %d", len(diags))
	}
}

func TestQSR011_MoreOption(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"test1.container",
		"[Container]\nImage=mock1\nPublishPort=42069:8080/tcp",
	)

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
		}
		diags := qsr011(s)

		if len(diags) != 0 {
			t.Fatalf("Exptected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}
