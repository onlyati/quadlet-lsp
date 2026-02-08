package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMacros_Valid tests when user type "new.A" then all possible
// option are popup.
func TestNewMacros_Valid(t *testing.T) {
	expected := map[string]bool{
		"new.Annotation": false,
		"new.AddHost":    false,
	}

	s := NewCompletion(
		[]string{"[Container]", "new.A"},
		"test.container",
		1,
		uint32(4),
	)
	s.config = &utils.QuadletConfig{}
	s.config.Podman = utils.BuildPodmanVersion(5, 5, 2)

	comps := listNewMacros(s)

	require.Len(t, comps, 2, "expected 2 completions")

	for _, d := range comps {
		if _, ok := expected[d.Label]; ok {
			expected[d.Label] = true
		} else {
			assert.Fail(t, "unexpected suggestion %s", d.Label)
		}
	}

	for k, v := range expected {
		if !v {
			assert.Fail(t, "did not get suggestion %s", k)
		}
	}
}
