package hover

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSystemdSpecifier tests hover on systemd identifiers.
func TestSystemdSpecifier(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "Volume=%h:%h:ro",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "Label=Test=Value%%%a",
			URI:               "file://test.volume",
			Section:           "Container",
			CharacterPosition: uint32(18),
		},
	}

	for _, info := range cases {
		hoverValue := handleSystemSpecifier(info)
		assert.NotNil(t, hoverValue, "expected hover value")
	}
}

// TestSystemdSpecifierEscaping tests that no hover when it is an escaped
// systemd identifier.
func TestSystemdSpecifierEscaping(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "Label=Test=%%a",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(12),
		},
		{
			Line:              "Label=Test=%%%%a",
			URI:               "file://test.volume",
			Section:           "Container",
			CharacterPosition: uint32(14),
		},
		{
			Line:              "%%%%a",
			URI:               "file://test.pod",
			Section:           "Container",
			CharacterPosition: uint32(3),
		},
	}

	for _, info := range cases {
		hoverValue := handleSystemSpecifier(info)
		assert.Nil(t, hoverValue, "expected hover value")
	}
}
