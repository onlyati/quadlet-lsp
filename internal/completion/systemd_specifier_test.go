package completion

import "testing"

func TestIsItSystemdSpecifier_Valid(t *testing.T) {
	cases := []struct {
		line    string
		charPos uint32
	}{
		{
			line:    "Label=\"systemd.unit=%\"",
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
		if !isItSystemSpecifier(s.line, s.charPos) {
			t.Fatalf("expected on systemd specifier at '%s'", s.line)
		}
	}
}
