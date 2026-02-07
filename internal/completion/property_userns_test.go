package completion

import (
	"os"
	"slices"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

type usernsMockCommander struct{}

func (m usernsMockCommander) Run(name string, args ...string) ([]string, error) {
	if args[2] == "scr.io/org/mock1:latest" {
		return []string{
			"[",
			"	{",
			"		 \"Config\": {",
			"			\"User\": \"999\" ",
			"		 }",
			"	}",
			"]",
		}, nil
	}

	return []string{}, nil
}

func TestPropertyUserIDs_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)
	s := NewCompletion(
		[]string{"[Container]", "UserNS=keep-id:", "Image=scr.io/org/mock1:latest"},
		"foo.container",
		1,
		0,
	)
	s.commander = usernsMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListUserIDs(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	if len(labels) != 2 {
		t.Fatalf("exptected 2 completions, but got %d", len(labels))
	}

	checkGID := slices.Contains(labels, "gid=999")
	checkUID := slices.Contains(labels, "uid=999")
	if !checkGID || !checkUID {
		t.Fatalf(
			"did not read correct values: %v %v %v",
			labels,
			checkGID,
			checkUID,
		)
	}
}

func TestPropertyUserIDs_Invalid(t *testing.T) {
	tmpDir := os.TempDir()
	_ = os.Chdir(tmpDir)
	s := NewCompletion(
		[]string{"[Container]", "UserNS=auto", "Image=scr.io/org/mock1:latest"},
		"foo.container",
		1,
		0,
	)
	s.commander = usernsMockCommander{}
	s.config = &utils.QuadletConfig{
		WorkspaceRoot: tmpDir,
		Project: utils.ProjectProperty{
			DirLevel: utils.ReturnAsPtr(2),
		},
	}

	comps := propertyListUserIDs(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	if len(labels) != 0 {
		t.Fatalf("exptected 0 completions, but got %d", len(labels))
	}
}
