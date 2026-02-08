package completion

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type secretMockCommander struct{}

func (c secretMockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{"secret1", "secret2"}, nil
}

// TestPropertySecret_SecretSuggestion tests if completion return with the existing
// secrets.
func TestPropertySecret_SecretSuggestion(t *testing.T) {
	s := Completion{}
	s.commander = secretMockCommander{}
	s.text = []string{"Secret="}
	s.line = 0
	s.char = uint32(len(s.text[0]))

	comps := propertyListSecrets(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 2, "expected length 2")
	assert.ElementsMatch(t, labels, []string{"secret1", "secret2"}, "did not read command output")
}

// TestPropertySecret_ColonSuggestion tests parameter completion of the secret.
func TestPropertySecret_ColonSuggestion(t *testing.T) {
	s := Completion{}
	s.commander = secretMockCommander{}
	s.text = []string{"Secret=secret1,"}
	s.line = 0
	s.char = uint32(len(s.text[0]))

	comps := propertyListSecrets(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 3, "expected length 3")
	assert.ElementsMatch(t, labels, []string{"type=mount", "type=env", "target="}, "did not show parameters")
}

// TestPropertySecret_CheckCursorPosition tests if secret completion is displayed
// based on the cursor position. So it is before the first ',' and after '='.
func TestPropertySecret_CheckCursorPosition(t *testing.T) {
	s := Completion{}
	s.commander = secretMockCommander{}
	s.text = []string{"Secret=,target=env,type=FOO"}
	s.line = 0
	s.char = uint32(strings.Index(s.text[0], ","))

	comps := propertyListSecrets(s)

	labels := []string{}
	for _, c := range comps {
		labels = append(labels, c.Label)
	}

	require.Len(t, labels, 2, "expected length 2")
	assert.ElementsMatch(t, labels, []string{"secret1", "secret2"}, "did not read command output")
}
