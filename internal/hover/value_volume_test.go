package hover

import (
	"testing"
)

func TestValueVolumeSource(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(8),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(20),
		},
	}

	for i, info := range cases {
		hoverValue := HoverFunction(info)

		if hoverValue == nil {
			t.Fatalf("expected hover value but got nil at #%d", i)
		}

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		if highlight != "/home/ati/tmp" {
			t.Fatalf("unexpected highlight but got '%s'", highlight)
		}
	}
}

func TestValueVolumeContainer(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(22),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(23),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			Uri:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(29),
		},
	}

	for i, info := range cases {
		hoverValue := HoverFunction(info)

		if hoverValue == nil {
			t.Fatalf("expected hover value but got nil at #%d", i)
		}

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		if highlight != "/app/tmp" {
			t.Fatalf("unexpected highlight but got '%s'", highlight)
		}
	}
}

func TestValueVolumeFlag(t *testing.T) {
	info := HoverInformation{
		Line:    "Volume=/home/ati/tmp:/app/tmp:rw,z,U,nocopy,shared",
		Uri:     "file://test.container",
		Section: "Container",
	}
	cases := []struct {
		position uint32
		flag     string
	}{
		{uint32(30), "rw"},
		{uint32(31), "rw"},
		{uint32(33), "z"},
		{uint32(35), "U"},
		{uint32(37), "nocopy"},
		{uint32(45), "shared"},
	}

	for i, c := range cases {
		info.CharacterPosition = c.position
		hoverValue := HoverFunction(info)

		if hoverValue == nil {
			t.Fatalf("expected hover value but got nil at #%d", i)
		}

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		if highlight != c.flag {
			t.Fatalf("unexpected highlight but got '%s'", highlight)
		}
	}
}
