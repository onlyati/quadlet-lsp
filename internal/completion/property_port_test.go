package completion

import (
	"os"
	"slices"
	"testing"
)

type portMockCommander struct{}

func (m portMockCommander) Run(name string, args ...string) ([]string, error) {
	if args[2] == "scr.io/org/mock1:latest" {
		return []string{
			"[",
			"	{",
			"		 \"Config\": {",
			"			\"ExposedPorts\": {",
			"				\"420/tcp\": {}",
			"			}",
			"		 }",
			"	}",
			"]",
		}, nil
	}
	if args[2] == "scr.io/org/mock2:latest" {
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

func TestPropertyPort_ValidRawImage(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"foo.container",
		"[Container]\nImage=scr.io/org/mock1:latest",
	)

	s := Completion{}
	s.commander = portMockCommander{}
	s.text = []string{"[Container]", "Image=scr.io/org/mock1:latest", "PublishPort=69:"}
	s.char = 0
	s.line = 2
	s.uri = "file://" + tmpDir + "/foo.container"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	if len(labels) != 1 {
		t.Fatalf("expected 1, but got %d", len(labels))
	}

	if labels[0] != "420" {
		t.Fatalf("exptected port 420, but got %s", labels[0])
	}
}

func TestPropertyPort_ValidImageFile(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"bar.container",
		"[Container]\nImage=bar.image",
	)
	createTempFile(
		t,
		tmpDir,
		"bar.image",
		"[Image]\nImage=scr.io/org/mock2:latest",
	)

	s := Completion{}
	s.commander = portMockCommander{}
	s.text = []string{"[Container]", "Image=bar.image", "PublishPort=69:"}
	s.char = 0
	s.line = 2
	s.uri = "file://" + tmpDir + "/bar.container"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	if len(labels) != 1 {
		t.Fatalf("expected 1, but got %d", len(labels))
	}

	if labels[0] != "69" {
		t.Fatalf("exptected port 69, but got %s", labels[0])
	}
}

func TestPropertyPort_ValidPod(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"bar.container",
		"[Container]\nPod=foobar.pod\nImage=bar.image",
	)
	createTempFile(
		t,
		tmpDir,
		"foo.container",
		"[Container]\nPod=foobar.pod\nImage=scr.io/org/mock1:latest",
	)
	createTempFile(
		t,
		tmpDir,
		"bar.image",
		"[Image]\nImage=scr.io/org/mock2:latest",
	)
	createTempFile(
		t,
		tmpDir,
		"foobar.pod",
		"[Pod]\nPublishPort=69:")

	s := Completion{}
	s.commander = portMockCommander{}
	s.text = []string{"[Pod]", "PublishPort=69:"}
	s.char = 0
	s.line = 1
	s.uri = "file://" + tmpDir + "/foobar.pod"

	comps := propertyListPorts(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	if len(labels) != 2 {
		t.Fatalf("expected 2, but got %d", len(labels))
	}

	checkPort69 := slices.Contains(labels, "69")
	checkPort420 := slices.Contains(labels, "420")
	if !checkPort69 || !checkPort420 {
		t.Fatalf(
			"Unexpected ports: %v %v %v",
			labels,
			checkPort420,
			checkPort69,
		)
	}
}
