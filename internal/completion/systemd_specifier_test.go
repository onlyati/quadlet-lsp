package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsItSystemdSpecifier_Valid tests completion for systemd identifier
// like %h, %t, and so on.
func TestIsItSystemdSpecifier_Valid(t *testing.T) {
	cases := []struct {
		line    string
		charPos uint32
	}{
		{
			line:    `Label="systemd.unit=%"`,
			charPos: 21,
		},
		{
			line:    "Label=systemd.unit=%",
			charPos: 20,
		},
		{
			line:    "Label=systemd.unit=%-unit",
			charPos: 20,
		},
	}

	for _, s := range cases {
		assert.True(t, isItSystemSpecifier(s.line, s.charPos), "expected on systemd specifier")
	}
}
