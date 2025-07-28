package completion

import (
	"slices"
	"strings"
	"testing"
)

type secretMockCommander struct{}

func (c secretMockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{"secret1", "secret2"}, nil
}

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

	containerSecret1 := slices.Contains(labels, "secret1")
	containerSecret2 := slices.Contains(labels, "secret2")
	if !containerSecret1 || !containerSecret2 {
		t.Fatalf("did not read command output: %v", labels)
	}
}

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

	containsType1Check := slices.Contains(labels, "type=mount")
	containsType2Check := slices.Contains(labels, "type=env")
	containerTargetCheck := slices.Contains(labels, "target=")
	if !containsType1Check || !containsType2Check || !containerTargetCheck {
		t.Fatalf(
			"did not select proper suggestions, %v %v %v",
			containsType1Check,
			containsType2Check,
			containerTargetCheck,
		)
	}
}

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

	containerSecret1 := slices.Contains(labels, "secret1")
	containerSecret2 := slices.Contains(labels, "secret2")
	if !containerSecret1 || !containerSecret2 {
		t.Fatalf("did not read command output: %v", labels)
	}
}
