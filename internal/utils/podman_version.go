package utils

import (
	"errors"
	"strconv"
	"strings"
)

type PodmanVersion struct {
	Version int
	Release int
	Minor   int
}

func NewPodmanVersion(c Commander) (PodmanVersion, error) {
	output, err := c.Run("podman", "version")
	if err != nil {
		return PodmanVersion{}, err
	}

	var rawVersion string
	for _, line := range output {
		if strings.HasPrefix(line, "Version:") {
			rawVersion = strings.TrimSpace(strings.Split(line, ":")[1])
			break
		}
	}

	return ParseVersion(rawVersion)
}

func ParseVersion(raw string) (PodmanVersion, error) {
	tmp := strings.Split(raw, ".")

	if len(tmp) != 3 {
		return PodmanVersion{}, errors.New("invalid version number")
	}

	version, err := strconv.Atoi(tmp[0])
	if err != nil {
		return PodmanVersion{}, err
	}

	release, err := strconv.Atoi(tmp[1])
	if err != nil {
		return PodmanVersion{}, err
	}

	minor, err := strconv.Atoi(tmp[2])
	if err != nil {
		return PodmanVersion{}, err
	}

	return PodmanVersion{
		Version: version,
		Release: release,
		Minor:   minor,
	}, nil
}
