package completion

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_ListNewProperties(t *testing.T) {
	cases := []Completion{
		{
			line: 0,
			text: []string{"H"},
			char: 1,
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 6, 0),
			},
			section: "Container",
		},
		{
			line: 0,
			text: []string{"Foo=bar", "H"},
			char: 1,
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 6, 0),
			},
			section: "Container",
		},
	}

	for _, s := range cases {
		result := listNewProperties(s)
		assert.Greater(t, len(result), 0)
	}
}

func Test_ListNewPropertiesWithCont(t *testing.T) {
	cases := []Completion{
		{
			line: 1,
			text: []string{"Foo=bar \\", "H"},
			char: 1,
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 6, 0),
			},
			section: "Container",
		},
	}

	for _, s := range cases {
		result := listNewProperties(s)
		assert.Equal(t, 0, len(result))
	}
}
