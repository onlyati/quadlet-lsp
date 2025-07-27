package utils

import (
	"encoding/json"
	"os"
	"path"
)

type QuadletConfig struct {
	Disable       []string      `json:"disable"`
	PodmanVersion string        `json:"podmanVersion"`
	Podman        PodmanVersion `json:"-"`
}

func LoadConfig(workspaceRoot string) (QuadletConfig, error) {
	configPath := path.Join(workspaceRoot, ".quadletrc.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return QuadletConfig{}, err
	}

	var config QuadletConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return QuadletConfig{}, err
	}

	if config.Podman, err = ParseVersion(config.PodmanVersion); err != nil {
		return QuadletConfig{}, err
	}

	return config, nil
}
