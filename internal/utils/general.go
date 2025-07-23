package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// I did not find better solution, probably not nice but works
func ReturnAsStringPtr(s string) *string {
	return &s
}

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

// List quadlet files from the current work directory based on extenstion
func ListQuadletFiles(ext string) ([]protocol.CompletionItem, error) {
	dirs := []protocol.CompletionItem{}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cwd = filepath.Join(cwd, ext)
	files, err := filepath.Glob(cwd)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		chunks := strings.Split(file, string(os.PathSeparator))
		dirs = append(dirs, protocol.CompletionItem{
			Label:         chunks[len(chunks)-1],
			Documentation: "From work directory: " + cwd,
		})
	}

	return dirs, nil
}
