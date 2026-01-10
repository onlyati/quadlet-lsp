package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
	"sync"
)

type QuadletConfig struct {
	Mu            sync.RWMutex    `json:"-"`
	Disable       []string        `json:"disable"`
	PodmanVersion string          `json:"podmanVersion"`
	Project       ProjectProperty `json:"project"`
	Podman        PodmanVersion   `json:"-"`
	WorkspaceRoot string          `json:"-"`
}

type ProjectProperty struct {
	DirLevel *int   `json:"dirLevel"`
	RootDir  string `json:"rootDir"`
}

func LoadConfig(workspaceRoot string, c Commander) (*QuadletConfig, error) {
	configPath := workspaceRoot
	if !strings.HasSuffix(workspaceRoot, ".quadletrc.json") {
		configPath = path.Join(workspaceRoot, ".quadletrc.json")
	}

	var config QuadletConfig
	data, err := os.ReadFile(configPath)
	if err == nil {
		err = json.Unmarshal(data, &config)
		if err != nil {
			return nil, err
		}
	}

	config.WorkspaceRoot = workspaceRoot
	if config.Project.RootDir != "" {
		config.WorkspaceRoot = path.Join(config.WorkspaceRoot, config.Project.RootDir)
	}

	if config.Podman, err = ParseVersion(config.PodmanVersion); err != nil {
		// try to discover podman version from the machine
		pVersion, err := NewPodmanVersion(c)
		if err != nil {
			config.Podman = BuildPodmanVersion(5, 4, 0)
		} else {
			config.Podman = pVersion
		}
	}

	// Check project properties
	if config.Project.RootDir == "" {
		config.Project.RootDir = config.WorkspaceRoot
	}
	if config.Project.DirLevel == nil {
		config.Project.DirLevel = ReturnAsPtr(2)
	}

	if *config.Project.DirLevel < 0 {
		return &config, errors.New("config.project.dirLevel must be greater than 0")
	}

	return &config, nil
}
