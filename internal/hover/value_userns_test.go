package hover

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValueUserNS tests hover output for the 'UserNS=' lines.
func TestValueUserNS(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "UserNS=keep-id",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "UserNS=keep-id",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(13),
		},
		{
			Line:              "UserNS=auto",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "UserNS=nomap",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "UserNS=host",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
	}

	for i, info := range cases {
		hoverValue := HoverFunction(info)

		require.NotNilf(t, hoverValue, "expected hover value at %d", i)

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		assert.Contains(t, []string{"keep-id", "nomap", "host", "auto"}, highlight, "unexpected highlight")
	}
}
