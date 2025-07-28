package completion

import (
	"os"
	"slices"
	"testing"
)

func TestPropertyPod_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.pod", "[Pod]")
	createTempFile(t, tmpDir, "bar.pod", "[Pod]")
	createTempFile(t, tmpDir, "foo.network", "[Network]")

	s := Completion{}

	comps := propertyListPods(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	checkFooNetwork := slices.Contains(labels, "foo.network")
	if checkFooNetwork {
		t.Fatalf("listed network but it should not: %v", labels)
	}

	checkFooPod := slices.Contains(labels, "foo.pod")
	checkBarPod := slices.Contains(labels, "bar.pod")
	if !checkFooPod || !checkBarPod {
		t.Fatalf(
			"did not list everything: %v %v %v",
			labels,
			checkFooPod,
			checkBarPod,
		)
	}
}
