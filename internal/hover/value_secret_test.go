package hover

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValueSecret tests the hover values on lines like 'Secret='.
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
				URI:     "file:///test.container",
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
				URI:     "file:///test.container",
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
				URI:     "file:///test.container",
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
				URI:     "file:///test.container",
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
				URI:     "file:///test.container",
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
				URI:     "file:///test.container",
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

		require.NotNilf(t, hoverValue, "expected hover value at %d", i)

		highlight := c.info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		assert.Equal(t, c.expectedHighlight.text, highlight, "unexpected highlight")
	}
}
