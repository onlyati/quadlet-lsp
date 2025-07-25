package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Commander interface {
	Run(name string, args ...string) ([]string, error)
}

type CommandExecutor struct{}

// This method execute an OS command and return with its output
func (c CommandExecutor) Run(name string, args ...string) ([]string, error) {
	output := []string{}
	cmd := exec.Command(name, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute podman: %w", err)
	}

	// Split output by newlines and filter out empty lines
	for line := range strings.SplitSeq(out.String(), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			output = append(output, trimmed)
		}
	}
	return output, nil
}
