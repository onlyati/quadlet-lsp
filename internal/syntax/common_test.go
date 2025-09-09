package syntax

import "testing"

func TestCanFileBeApplied_Valid(t *testing.T) {
	list := []string{"container", "volume", "image"}

	tmp := canFileBeApplied("file:///foo.container", list)
	if tmp != "[Container]" {
		t.Fatalf("expected 'Container', got '%s'", tmp)
	}

	tmp = canFileBeApplied("file:///foo.container.d/10-ports.conf", list)
	if tmp != "[Container]" {
		t.Fatalf("expected 'Container', got '%s'", tmp)
	}

	tmp = canFileBeApplied("file:///foo.image", list)
	if tmp != "[Image]" {
		t.Fatalf("expected 'Image', got '%s'", tmp)
	}

	tmp = canFileBeApplied("file:///foo.image.d/auth.conf", list)
	if tmp != "[Image]" {
		t.Fatalf("expected 'Image', got '%s'", tmp)
	}

	tmp = canFileBeApplied("file:///foo.volume", list)
	if tmp != "[Volume]" {
		t.Fatalf("expected 'Volume', got '%s'", tmp)
	}

	tmp = canFileBeApplied("file:///foo.volume.d/driver.conf", list)
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

	tmp = canFileBeApplied("file:///foo.conf", list)
	if tmp != "" {
		t.Fatalf("expected '', got '%s'", tmp)
	}
}
