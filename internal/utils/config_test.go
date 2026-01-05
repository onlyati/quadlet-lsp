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
	assert.Equal(t, tmpDir, cfg.WorkspaceRoot, "workspace root should be the current directory")
	assert.Equal(t, utils.BuildPodmanVersion(5, 5, 2), cfg.Podman, "wrong podman version gathered")
	assert.Nil(t, cfg.Disable, "cfg.disable should be nil")
	assert.Equal(t, cfg.WorkspaceRoot, cfg.Project.RootDir, "if no project rootDir specified, this should be workspace root")
	assert.Equal(t, 2, *cfg.Project.DirLevel, "dirLevel default is 2")
}

func TestConfig_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	createTempFile(t, tmpDir, ".quadletrc.json", "{ \"podmanVersion\": \"5.4.0\", \"disable\": [\"qsr003\"] }")

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.NoError(t, err)
	assert.Equal(t, tmpDir, cfg.WorkspaceRoot, "workspace root should be the current directory")
	assert.Equal(t, utils.BuildPodmanVersion(5, 4, 0), cfg.Podman, "wrong podman version gathered")
	assert.Len(t, cfg.Disable, 1, "wrong disable array length")
	assert.Equal(t, "qsr003", cfg.Disable[0], "expected qsr003 rule")
	assert.Equal(t, cfg.WorkspaceRoot, cfg.Project.RootDir, "if no project rootDir specified, this should be workspace root")
	assert.Equal(t, 2, *cfg.Project.DirLevel, "dirLevel default is 2")
}

func TestConfig_WithProjectProps(t *testing.T) {
	tmpDir := t.TempDir()
	createTempFile(t, tmpDir, ".quadletrc.json", `{ "project" : { "rootDir": "./containers", "dirLevel": 4 } }`)

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.NoError(t, err)
	assert.Equal(t, tmpDir, cfg.WorkspaceRoot, "workspace root should be the current directory")
	assert.Equal(t, utils.BuildPodmanVersion(5, 5, 2), cfg.Podman, "wrong podman version gathered")
	assert.Nil(t, cfg.Disable, "wrong disable array length")
	assert.Equal(t, "./containers", cfg.Project.RootDir, "if no project rootDir specified, this should be workspace root")
	assert.Equal(t, 4, *cfg.Project.DirLevel, "dirLevel default is 2")
}

func TestConfig_InvalidJson(t *testing.T) {
	tmpDir := t.TempDir()
	createTempFile(t, tmpDir, ".quadletrc.json", "{ \"podmanVersion\": \"5.4.0\", \"disable\": [\"qsr003\"] ")

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.Error(t, err)
	assert.Nil(t, cfg, "config should be nil")
}

func TestConfig_PodmanNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	c := configNotFoundMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.NoError(t, err)
	assert.NotNil(t, cfg, "config should be nil")
	assert.Equal(t, tmpDir, cfg.WorkspaceRoot, "workspace root should be the current directory")
	assert.Equal(t, utils.BuildPodmanVersion(5, 4, 0), cfg.Podman, "If podman not found then 5.4.0 is the default")
	assert.Nil(t, cfg.Disable, "cfg.disable should be nil")
	assert.Equal(t, cfg.WorkspaceRoot, cfg.Project.RootDir, "if no project rootDir specified, this should be workspace root")
	assert.Equal(t, 2, *cfg.Project.DirLevel, "dirLevel default is 2")
}

func TestConfig_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	c := configMockCommander{}
	cfg, err := utils.LoadConfig(tmpDir, c)
	assert.NoError(t, err)
	assert.NotNil(t, cfg, "config should be nil")
	assert.Equal(t, tmpDir, cfg.WorkspaceRoot, "workspace root should be the current directory")
	assert.Equal(t, utils.BuildPodmanVersion(5, 5, 2), cfg.Podman, "5.5.2 should be read from mock")
	assert.Nil(t, cfg.Disable, "cfg.disable should be nil")
	assert.Equal(t, cfg.WorkspaceRoot, cfg.Project.RootDir, "if no project rootDir specified, this should be workspace root")
	assert.Equal(t, 2, *cfg.Project.DirLevel, "dirLevel default is 2")
}
