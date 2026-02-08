package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCanFileBeApplied_Valid tests the extension checker function.
func TestCanFileBeApplied_Valid(t *testing.T) {
	list := []string{"container", "volume", "image"}

	assert.Equal(t, "[Container]", canFileBeApplied("file:///foo.container", list))
	assert.Equal(t, "[Container]", canFileBeApplied("file:///foo.container.d/10-ports.conf", list))
	assert.Equal(t, "[Image]", canFileBeApplied("file:///foo.image", list))
	assert.Equal(t, "[Image]", canFileBeApplied("file:///foo.image.d/auth.conf", list))
	assert.Equal(t, "[Volume]", canFileBeApplied("file:///foo.volume", list))
	assert.Equal(t, "[Volume]", canFileBeApplied("file:///foo.volume.d/driver.conf", list))
}

// TestCanFileBeApplied_Invalid tests when extension is forbidden.
func TestCanFileBeApplied_Invalid(t *testing.T) {
	list := []string{"container", "volume", "image"}

	assert.Equal(t, "", canFileBeApplied("file:///foo.network", list))
	assert.Equal(t, "", canFileBeApplied("file:///foo.conf", list))
}
