package parser

import (
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
