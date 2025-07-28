package completion

import (
	"testing"
)

func TestSection_Valid(t *testing.T) {
	expected := []string{
		"Pod",
		"Kube",
		"Network",
		"Volume",
		"Image",
		"Container",
	}

	diags := listSections(Completion{})
	result := []string{}

	for _, diag := range diags {
		result = append(result, diag.Label)
	}

	if len(expected) != len(result) {
		t.Fatalf("Expected '%v', but got '%v'", expected, result)
	}
}
