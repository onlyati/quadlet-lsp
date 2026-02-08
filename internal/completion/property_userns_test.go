package completion

import (
	"os"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type usernsMockCommander struct{}

func (m usernsMockCommander) Run(name string, args ...string) ([]string, error) {
	if args[2] == "scr.io/org/mock1:latest" {
		return []string{
			`[`,
			`	{`,
			`		 "Config": {`,
			`			"User": "999" `,
			`		 }`,
			`	}`,
			`]`,
		}, nil
	}

	return []string{}, nil
}

func TestPropertyUserIDs_Valid(t *testing.T) {
	tmpDir := t.TempDir()

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

	require.Len(t, labels, 2, "expected length 2")
	assert.ElementsMatch(t, labels, []string{"uid=999", "gid=999"}, "unexpected userns parameters")
}

// TestPropertyUserIDs_Invalid tests if no completion on case of invalid userns.
func TestPropertyUserIDs_Invalid(t *testing.T) {
	tmpDir := os.TempDir()

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

	require.Len(t, labels, 0, "should not be any completion")
}
