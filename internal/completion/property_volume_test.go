package completion

import (
	"os"
	"slices"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

type volumeMockCommnander struct{}

func (c volumeMockCommnander) Run(name string, args ...string) ([]string, error) {
	return []string{"volume1", "volume2"}, nil
}

func TestPropertyVolume_ListVolume(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.volume", "[Volume]")
	createTempFile(t, tmpDir, "foo.network", "[Network]")

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

	checkFooNetwork := slices.Contains(labels, "foo.network")
	if checkFooNetwork {
		t.Fatalf("listed network but it should not: %v", labels)
	}

	checkFooVolume := slices.Contains(labels, "foo.volume")
	checkVolume1 := slices.Contains(labels, "volume1")
	checkVolume2 := slices.Contains(labels, "volume2")
	if !checkFooVolume || !checkVolume1 || !checkVolume2 {
		t.Fatalf(
			"did not list everything: %v %v %v %v",
			labels,
			checkFooVolume,
			checkVolume1,
			checkVolume2,
		)
	}
}

func TestPropertyVolume_NoList(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.volume", "[Volume]")
	createTempFile(t, tmpDir, "foo.network", "[Network]")

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

	if len(comps) > 0 {
		t.Fatalf("expected 0 completions, but got %d", len(comps))
	}
}

func TestPropertyVolume_ListFlags(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.volume", "[Volume]")
	createTempFile(t, tmpDir, "foo.network", "[Network]")

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

	checkLabelRw := slices.Contains(labels, "rw")
	checkLabelRo := slices.Contains(labels, "ro")
	checkLabelZz := slices.Contains(labels, "z")
	checkLabelZZ := slices.Contains(labels, "Z")
	if !checkLabelRw || !checkLabelRo || !checkLabelZZ || !checkLabelZz {
		t.Fatalf(
			"Unexpected suggestions returned: %v %v %v %v %v",
			labels,
			checkLabelRw,
			checkLabelRo,
			checkLabelZz,
			checkLabelZZ,
		)
	}
}

func TestPropertyVolume_ListVolumeWithCursor(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.volume", "[Volume]")
	createTempFile(t, tmpDir, "foo.network", "[Network]")

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

	checkFooNetwork := slices.Contains(labels, "foo.network")
	if checkFooNetwork {
		t.Fatalf("listed network but it should not: %v", labels)
	}

	checkFooVolume := slices.Contains(labels, "foo.volume")
	checkVolume1 := slices.Contains(labels, "volume1")
	checkVolume2 := slices.Contains(labels, "volume2")
	if !checkFooVolume || !checkVolume1 || !checkVolume2 {
		t.Fatalf(
			"did not list everything: %v %v %v %v",
			labels,
			checkFooVolume,
			checkVolume1,
			checkVolume2,
		)
	}
}
