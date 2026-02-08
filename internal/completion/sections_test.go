package completion

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSection_Valid tests if sections list is valid.
func TestSection_Valid(t *testing.T) {
	expected := []string{
		"Pod",
		"Kube",
		"Build",
		"Network",
		"Volume",
		"Image",
		"Container",
		"Artifact",
		"Unit",
		"Install",
	}

	diags := listSections(Completion{})
	result := []string{}

	for _, diag := range diags {
		result = append(result, diag.Label)
	}

	require.Len(t, result, len(expected), "did not contain all sections")
}
