package hover

import (
	"testing"
)

func TestValueSecret(t *testing.T) {
	cases := []struct {
		info              HoverInformation
		expectedHighlight struct {
			position uint32
			text     string
		}
	}{
		{
			info: HoverInformation{
				Line:    "Secret=my-secret,type=env,target=MYTARGET",
				Uri:     "file:///test.container",
				Section: "Container",
			},
			expectedHighlight: struct {
				position uint32
				text     string
			}{
				7, "my-secret",
			},
		},
		{
			info: HoverInformation{
				Line:    "Secret=my-secret,type=env,target=MYTARGET",
				Uri:     "file:///test.container",
				Section: "Container",
			},
			expectedHighlight: struct {
				position uint32
				text     string
			}{
				15, "my-secret",
			},
		},
		{
			info: HoverInformation{
				Line:    "Secret=my-secret,type=env,target=MYTARGET",
				Uri:     "file:///test.container",
				Section: "Container",
			},
			expectedHighlight: struct {
				position uint32
				text     string
			}{
				17, "type=env",
			},
		},
		{
			info: HoverInformation{
				Line:    "Secret=my-secret,type=env,target=MYTARGET",
				Uri:     "file:///test.container",
				Section: "Container",
			},
			expectedHighlight: struct {
				position uint32
				text     string
			}{
				24, "type=env",
			},
		},
		{
			info: HoverInformation{
				Line:    "Secret=my-secret,type=env,target=MYTARGET",
				Uri:     "file:///test.container",
				Section: "Container",
			},
			expectedHighlight: struct {
				position uint32
				text     string
			}{
				26, "target=MYTARGET",
			},
		},
		{
			info: HoverInformation{
				Line:    "Secret=my-secret,type=env,target=MYTARGET",
				Uri:     "file:///test.container",
				Section: "Container",
			},
			expectedHighlight: struct {
				position uint32
				text     string
			}{
				40, "target=MYTARGET",
			},
		},
	}

	for i, c := range cases {
		c.info.CharacterPosition = c.expectedHighlight.position
		hoverValue := HoverFunction(c.info)

		if hoverValue == nil {
			t.Fatalf("expected hover value but got nil at #%d", i)
		}

		isItStart := hoverValue.Range.Start.Character == c.expectedHighlight.position
		isItEnd := hoverValue.Range.End.Character-1 == c.expectedHighlight.position
		if !isItEnd && !isItStart {
			t.Fatalf(
				"charatcer position should be at boundary but it %d id not %d-%d",
				c.expectedHighlight.position,
				hoverValue.Range.Start.Character,
				hoverValue.Range.End.Character,
			)
		}

		highlight := c.info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		if highlight != c.expectedHighlight.text {
			t.Fatalf("unexpected highlight but got '%s'", highlight)
		}
	}
}
