package lsp

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func execPodmanCommand(params []string) ([]string, error) {
	output := []string{}

	cmd := exec.Command("podman", params...)

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

func listQuadletFiles(ext string) ([]protocol.CompletionItem, error) {
	dirs := []protocol.CompletionItem{}

	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	return nil, err
	// }
	// homeDir = filepath.Join(homeDir, ".config", "containers", "systemd", ext)
	// files, err := filepath.Glob(homeDir)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, file := range files {
	// 	chunks := strings.Split(file, string(os.PathSeparator))
	// 	dirs = append(dirs, protocol.CompletionItem{
	// 		Label:         chunks[len(chunks)-1],
	// 		Documentation: "From home: " + homeDir,
	// 	})
	// }

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

func findSection(lines []string, lineNumber protocol.UInteger) string {
	section := ""
	for i := lineNumber; ; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			break
		}

		if i == 0 {
			break
		}
	}
	return section
}

func returnAsStringPtr(s string) *string {
	return &s
}
