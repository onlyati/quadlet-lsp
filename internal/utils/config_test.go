package utils_test

import (
	"errors"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
)

type configMockCommander struct{}

func (c configMockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{
		"Client:        Podman Engine",
		"Version:       5.5.2",
		"API Version:   5.5.2",
		"Go Version:    go1.24.4",
		"Git Commit:    e7d8226745ba07a64b7176a7f128e4ef53225a0e",
		"Built:         Tue Jun 24 02:00:00 2025",
		"Build Origin:  Fedora Project",
		"OS/Arch:       linux/amd64",
	}, nil
}

type configNotFoundMockCommander struct{}

func (c configNotFoundMockCommander) Run(name string, args ...string) ([]string, error) {
	return []string{}, errors.New("command not found")
}

func TestConfig_Default(t *testing.T) {
	tmpDir := t.TempDir()

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.NoError(t, err)

	if cfg.WorkspaceRoot != tmpDir {
		t.Fatalf("Expected WorkspaceRoot '%s', but got '%s'", tmpDir, cfg.WorkspaceRoot)
	}

	if cfg.Podman != utils.BuildPodmanVersion(5, 5, 2) {
		t.Fatalf("Expected Podman 5.5.2, but got: %+v", cfg.Podman)
	}

	if cfg.Disable != nil {
		t.Fatalf("Expected nil, but got %v", cfg.Disable)
	}
}

func TestConfig_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, ".quadletrc.json", "{ \"podmanVersion\": \"5.4.0\", \"disable\": [\"qsr003\"] }")

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.NoError(t, err)

	if cfg.WorkspaceRoot != tmpDir {
		t.Fatalf("Expected WorkspaceRoot '%s', but got '%s'", tmpDir, cfg.WorkspaceRoot)
	}

	if cfg.Podman != utils.BuildPodmanVersion(5, 4, 0) {
		t.Fatalf("Expected Podman 5.4.0, but got: %+v", cfg.Podman)
	}

	if len(cfg.Disable) != 1 {
		t.Fatalf("Expected 1 disabled rule but got %d", len(cfg.Disable))
	}

	if cfg.Disable[0] != "qsr003" {
		t.Fatalf("Expected qsr003 disable rule, but got %s", cfg.Disable[0])
	}
}

func TestConfig_InvalidJson(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, ".quadletrc.json", "{ \"podmanVersion\": \"5.4.0\", \"disable\": [\"qsr003\"] ")

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.Error(t, err)

	if cfg != nil {
		t.Fatalf("Expected a nil config but got: %+v", cfg)
	}
}

func TestConfig_PodmanNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	c := configNotFoundMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.Error(t, err)

	if cfg == nil {
		t.Fatal("Expected but got nil")
	}

	if cfg.Podman != utils.BuildPodmanVersion(5, 4, 0) {
		t.Fatalf("Expected Podman 5.4.0, but got: %+v", cfg.Podman)
	}
}

func TestCOnfig_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.NoError(t, err)

	if cfg == nil {
		t.Fatal("Expected but got nil")
	}

	if cfg.WorkspaceRoot != tmpDir {
		t.Fatalf("Expected WorkspaceRoot '%s', but got '%s'", tmpDir, cfg.WorkspaceRoot)
	}

	if cfg.Podman != utils.BuildPodmanVersion(5, 5, 2) {
		t.Fatalf("Expected Podman 5.5.2, but got: %+v", cfg.Podman)
	}

	if cfg.Disable != nil {
		t.Fatalf("Expected nil, but got %v", cfg.Disable)
	}
}
