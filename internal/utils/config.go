package utils

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"sync"
)

type QuadletConfig struct {
	Mu            sync.RWMutex  `json:"-"`
	Disable       []string      `json:"disable"`
	PodmanVersion string        `json:"podmanVersion"`
	Podman        PodmanVersion `json:"-"`
	WorkspaceRoot string        `json:"-"`
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
	if config.Podman, err = ParseVersion(config.PodmanVersion); err != nil {
		// try to discover podman version from the machine
		pVersion, err := NewPodmanVersion(c)
		if err != nil {
			config.Podman = BuildPodmanVersion(5, 4, 0)
			return &config, err
		}
		config.Podman = pVersion
	}

	return &config, nil
}
