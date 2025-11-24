package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_gatherSectionName(t *testing.T) {
	source := []string{
		"[Container]",
		"[Unit]",
		"[Install]",
	}
	expected := []string{
		"Container",
		"Unit",
		"Install",
	}

	for i, s := range source {
		result, found := gatherSectionName(s)
		assert.Equal(
			t,
			true,
			found,
			"exptects to found section",
		)
		assert.Equal(
			t,
			expected[i],
			result,
			"gathered section has unexpected value",
		)
	}
}

func Test_isDropinsBelongsToQuadlet(t *testing.T) {
	type testCase struct {
		parentDir     string
		possibleOwner string
		expected      bool
	}

	cases := []testCase{
		{"foo.container.d", "foo.container", true},
		{"container.d", "foo.container", true},
		{"foo-.container.d", "foo-bar-app.container", true},
		{"foo-bar-.container.d", "foo-bar-app.container", true},
		{"foo-bar-app.container.d", "foo-bar-app.container", true},
	}

	for _, c := range cases {
		result := isDropinsBelongsToQuadlet(c.possibleOwner, c.parentDir)
		assert.Equal(
			t, c.expected, result, fmt.Sprintf("case: %+v", c),
		)
	}
}
