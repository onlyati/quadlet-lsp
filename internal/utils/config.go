package utils

import (
	"encoding/json"
	"log"
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

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config QuadletConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Podman, err = ParseVersion(config.PodmanVersion); err != nil {
		log.Printf("Podman parse error: %v", err.Error())
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
