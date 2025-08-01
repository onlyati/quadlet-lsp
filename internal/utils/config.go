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

func LoadConfig(workspaceRoot string) (*QuadletConfig, error) {
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

	if config.Podman, err = ParseVersion(config.PodmanVersion); err != nil {
		// try to discover podman version from the machine
		c := CommandExecutor{}
		pVersion, err := NewPodmanVersion(c)
		if err != nil {
			return nil, err
		}
		config.Podman = pVersion
	}
	config.WorkspaceRoot = workspaceRoot

	return &config, nil
}
