package completion

import (
	"os"
	"slices"
	"testing"
)

type networkMockCommnander struct{}

func (c networkMockCommnander) Run(name string, args ...string) ([]string, error) {
	return []string{"network1", "network2"}, nil
}

func TestPropertyNetwork_ListNetwork(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.network", "[Network]")
	createTempFile(t, tmpDir, "foo.volume", "[Volume]")

	s := NewCompletion(
		[]string{"Network="},
		"test.container",
		0,
		uint32(len("Network=")),
	)
	s.commander = networkMockCommnander{}

	comps := propertyListNetworks(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	checkFooVolume := slices.Contains(labels, "foo.volume")
	if checkFooVolume {
		t.Fatalf("listed volume but it should not: %v", labels)
	}

	checkFooNetwork := slices.Contains(labels, "foo.network")
	checkNetwork1 := slices.Contains(labels, "network1")
	checkNetwork2 := slices.Contains(labels, "network2")
	if !checkFooNetwork || !checkNetwork1 || !checkNetwork2 {
		t.Fatalf(
			"did not list everything: %v %v %v %v",
			labels,
			checkFooNetwork,
			checkNetwork1,
			checkNetwork2,
		)
	}
}
