package utils_test

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

type mockCommander struct{}

func (c mockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{
		"Client:        Podman Engine",
		"Version:       5.4.2",
		"API Version:   5.4.2",
		"Go Version:    go1.24.4",
		"Git Commit:    e7d8226745ba07a64b7176a7f128e4ef53225a0e",
		"Built:         Tue Jun 24 02:00:00 2025",
		"Build Origin:  Fedora Project",
		"OS/Arch:       linux/amd64",
	}, nil
}

func TestNewPodmanVersion(t *testing.T) {
	expected := utils.PodmanVersion{
		Version: 5,
		Release: 4,
		Minor:   2,
	}

	result, err := utils.NewPodmanVersion(mockCommander{})
	assert.NoError(t, err)

	if expected != result {
		t.Fatalf("Exptected: '%v', but got '%v'", expected, result)
	}
}

func TestPodmanVersionGreateThan(t *testing.T) {
	p := utils.PodmanVersion{
		Version: 5,
		Release: 4,
		Minor:   2,
	}

	if p.GreaterOrEqual(utils.PodmanVersion{Version: 5, Release: 4, Minor: 0}) == false {
		t.Fatal("failed test case 1")
	}

	if p.GreaterOrEqual(utils.PodmanVersion{Version: 5, Release: 5, Minor: 0}) == true {
		t.Fatal("failed test case 2")
	}
}
