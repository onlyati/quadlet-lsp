package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

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

	if len(comps) < 2 {
		t.Fatalf("Exptected at least 2 completions, got %d", len(comps))
	}

	for _, d := range comps {
		if _, ok := expected[d.Label]; ok {
			expected[d.Label] = true
		} else {
			t.Fatalf("unexpected suggestion: %s", d.Label)
		}
	}

	for k, v := range expected {
		if !v {
			t.Fatalf("did not get suggestion: %s", k)
		}
	}
}
