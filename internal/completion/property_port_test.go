package completion

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type portMockCommander struct{}

func (m portMockCommander) Run(name string, args ...string) ([]string, error) {
	if args[2] == "scr.io/org/mock1:latest" {
		return []string{
			`[`,
			`	{`,
			`		 "Config": {`,
			`			"ExposedPorts": {`,
			`				"420/tcp": {}`,
			`			}`,
			`		 }`,
			`	}`,
			`]`,
		}, nil
	}
	if args[2] == "scr.io/org/mock2:latest" {
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

// TestPropertyPort_ValidRawImage test expects to return with ports that are
// read from the image using `podman image inspect` mock command.
func TestPropertyPort_ValidRawImage(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(
		t,
		tmpDir,
		"foo.container",
		"[Container]\nImage=scr.io/org/mock1:latest",
	)

	s := Completion{}
	s.commander = portMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	s.text = []string{"[Container]", "Image=scr.io/org/mock1:latest", "PublishPort=69:"}
	s.char = 14
	s.line = 2
	s.uri = "file://" + tmpDir + "/foo.container"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 1, "expected length 1")
	assert.Equal(t, labels[0], "420", "expected port 420")
}

// TestPropertyPort_ValidImageFile test read the port from the image file.
// Normally if the Image=bar.image, then it read the image file and check
// the image in that file.
func TestPropertyPort_ValidImageFile(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "bar.container", "[Container]\nImage=bar.image")
	testutils.CreateTempFile(t, tmpDir, "bar.image", "[Image]\nImage=scr.io/org/mock2:latest")

	s := Completion{}
	s.commander = portMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	s.text = []string{"[Container]", "Image=bar.image", "PublishPort=69:"}
	s.char = 0
	s.line = 2
	s.uri = "file://" + tmpDir + "/bar.container"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 1, "expected length 1")
	assert.Equal(t, labels[0], "69", "expected port 69")
}

// TestPropertyPort_ValidPod test is checking port from the pod file. To achive
// it, the program must looking for container units that belongs to this pod and
// check the image or image file that is that container.
func TestPropertyPort_ValidPod(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "bar.container", "[Container]\nPod=foobar.pod\nImage=bar.image")
	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nPod=foobar.pod\nImage=scr.io/org/mock1:latest")
	testutils.CreateTempFile(t, tmpDir, "bar.image", "[Image]\nImage=scr.io/org/mock2:latest")
	testutils.CreateTempFile(t, tmpDir, "foobar.pod", "[Pod]\nPublishPort=69:")

	s := Completion{}
	s.commander = portMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	s.text = []string{"[Pod]", "PublishPort=69:"}
	s.char = 0
	s.line = 1
	s.uri = "file://" + tmpDir + "/foobar.pod"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 2, "expected length 2")
	assert.ElementsMatch(t, labels, []string{"69", "420"}, "Unexpected ports")
}

// TestPropertyPort_ValidRawImageInDropins1 test is looking for the image in the
// container file meanwhile typing in drop-ins directory.
func TestPropertyPort_ValidRawImageInDropins1(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nImage=scr.io/org/mock1:latest")
	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.container.d"), "image.conf", "[Container]\nPublishPort=69:")

	s := Completion{}
	s.commander = portMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	s.text = []string{"[Container]", "PublishPort=69:"}
	s.char = 14
	s.line = 1
	s.uri = "file://" + tmpDir + "/foo.container.d/image.conf"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 1, "expected length 1")
	assert.Equal(t, labels[0], "420", "expected port 420")
}

// TestPropertyPort_ValidRawImageInDropins2 test is looking for the image in the
// drop-ins directory meanwhile typing there.
func TestPropertyPort_ValidRawImageInDropins2(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.container", "")
	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.container.d"), "image.conf", "[Container]\nImage=scr.io/org/mock1:latest")

	s := Completion{}
	s.commander = portMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}
	s.text = []string{"[Container]", "PublishPort=69:"}
	s.char = 14
	s.line = 1
	s.uri = "file://" + tmpDir + "/foo.container.d/image.conf"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 1, "expected length 1")
	assert.Equal(t, labels[0], "420", "expected port 420")
}
