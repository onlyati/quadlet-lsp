package hover

import (
	"testing"
)

func TestValueUserNS(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "UserNS=keep-id",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "UserNS=keep-id",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(13),
		},
		{
			Line:              "UserNS=auto",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "UserNS=nomap",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "UserNS=host",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
	}

	for i, info := range cases {
		hoverValue := HoverFunction(info)

		if hoverValue == nil {
			t.Fatalf("expected hover value but got nil at #%d", i)
		}

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		if highlight != "keep-id" && highlight != "nomap" && highlight != "host" && highlight != "auto" {
			t.Fatalf("unexpected highlight but got '%s'", highlight)
		}
	}
}
