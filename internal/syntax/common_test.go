package syntax

import "testing"

func TestCanFileBeApplied_Valid(t *testing.T) {
	list := []string{"container", "volume", "image"}

	tmp := canFileBeApplied("file:///foo.container", list)
	if tmp != "[Container]" {
		t.Fatalf("expected 'Container', got '%s'", tmp)
	}

	tmp = canFileBeApplied("file:///foo.image", list)
	if tmp != "[Image]" {
		t.Fatalf("expected 'Image', got '%s'", tmp)
	}

	tmp = canFileBeApplied("file:///foo.volume", list)
	if tmp != "[Volume]" {
		t.Fatalf("expected 'Volume', got '%s'", tmp)
	}
}

func TestCanFileBeApplied_Invalid(t *testing.T) {
	list := []string{"container", "volume", "image"}

	tmp := canFileBeApplied("file:///foo.network", list)
	if tmp != "" {
		t.Fatalf("expected '', got '%s'", tmp)
	}
}

func TestNamingRegexp(t *testing.T) {
	valid := []string{"foo", "bar"}
	for _, s := range valid {
		if !namingConvention.MatchString(s) {
			t.Fatalf("Regexp should match but it does not: %v %s", namingConvention, s)
		}
	}

	invalid := []string{".foo", "*bar", "_foo", "-bar"}
	for _, s := range invalid {
		if namingConvention.MatchString(s) {
			t.Fatalf("Regexp should not match but it does not: %v %s", namingConvention, s)
		}
	}
}
