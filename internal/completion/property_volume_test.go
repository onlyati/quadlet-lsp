package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type volumeMockCommnander struct{}

func (c volumeMockCommnander) Run(name string, args ...string) ([]string, error) {
	return []string{"volume1", "volume2"}, nil
}

// TestPropertyVolume_ListVolume tests if only *.volume files are displayed on
// completion of Volume.
func TestPropertyVolume_ListVolume(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]")
	testutils.CreateTempFile(t, tmpDir, "foo.network", "[Network]")

	s := NewCompletion(
		[]string{"Volume="},
		"test.container",
		0,
		uint32(len("Volume=")),
	)
	s.commander = volumeMockCommnander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListVolumes(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	assert.NotContains(t, labels, "foo.network", "should only list volmes")
	assert.ElementsMatch(t, labels, []string{"foo.volume", "volume1", "volume2"}, "did not list everything")
}

// TestPropertyVolume_NoList tests if no completion after the first ':' sign in
// the Volume line.
func TestPropertyVolume_NoList(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]")
	testutils.CreateTempFile(t, tmpDir, "foo.network", "[Network]")

	s := NewCompletion(
		[]string{"Volume=foo.volume:"},
		"test.container",
		0,
		uint32(len("Volume=foo.volume:")),
	)
	s.commander = volumeMockCommnander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListVolumes(s)

	require.Len(t, comps, 0, "expected 0 completion")
}

// TestPropertyVolume_ListFlags tests if flags are displayed after the second
// ':' sign.
func TestPropertyVolume_ListFlags(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]")
	testutils.CreateTempFile(t, tmpDir, "foo.network", "[Network]")

	s := NewCompletion(
		[]string{"Volume=foo.volume:/app/:"},
		"test.container",
		0,
		uint32(len("Volume=foo.volume:/app/:")),
	)
	s.commander = volumeMockCommnander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListVolumes(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Greater(t, len(labels), 0, "expected flags on completion")

	expectedFlags := []string{
		"rw",
		"ro",
		"z",
		"Z",
		"O",
		"copy",
		"nocopy",
		"dev",
		"nodev",
		"exec",
		"noexec",
		"suid",
		"nosuid",
		"bind",
		"rbind",
		"slave",
		"rslave",
		"shared",
		"rshared",
		"private",
		"rprivate",
		"unbindable",
		"runbindable",
	}
	assert.ElementsMatch(t, labels, expectedFlags, "unexpected suggestions")
}

// TestPropertyVolume_ListVolumeWithCursor tests if volume completions are displayed
// if cursor before the first ':' sign.
func TestPropertyVolume_ListVolumeWithCursor(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]")
	testutils.CreateTempFile(t, tmpDir, "foo.network", "[Network]")

	s := NewCompletion(
		[]string{"Volume=:/app/:"},
		"test.container",
		0,
		uint32(len("Volume=")),
	)
	s.commander = volumeMockCommnander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListVolumes(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	assert.NotContains(t, labels, "foo.network", "network should not be listed but volume")
	assert.ElementsMatch(t, labels, []string{"foo.volume", "volume1", "volume2"}, "did not list everything")
}
