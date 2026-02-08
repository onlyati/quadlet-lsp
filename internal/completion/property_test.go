package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_IsItPropertyCompletion tests condition of property completion call.
func Test_IsItPropertyCompletion(t *testing.T) {
	cases := []Completion{
		{
			line: 0,
			text: []string{"Foo="},
			char: 4,
		},
		{
			line: 0,
			text: []string{"Foo=bar"},
			char: 6,
		},
	}

	for _, s := range cases {
		result := isItPropertyCompletion(s)
		assert.True(t, result, "test should have been valid")
	}
}
