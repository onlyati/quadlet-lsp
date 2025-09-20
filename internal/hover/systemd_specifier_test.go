package hover

import "testing"

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

		if hoverValue == nil {
			t.Fatalf("expected hover value but got nil at %s", info.URI)
		}
	}
}

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

		if hoverValue != nil {
			t.Fatalf("not expected hover value but got %+v at %s", hoverValue, info.URI)
		}

	}
}
