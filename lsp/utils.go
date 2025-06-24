package lsp

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Commander interface {
	Run(name string, args ...string) ([]string, error)
}

type CommandExecutor struct{}

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

func listQuadletFiles(ext string) ([]protocol.CompletionItem, error) {
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

func findLineStartWith(prefix string) ([]protocol.Location, error) {
	var locations []protocol.Location

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, prefix) {
				locations = append(locations, protocol.Location{
					URI: protocol.DocumentUri("file://" + absPath),
					Range: protocol.Range{
						Start: protocol.Position{Line: protocol.UInteger(lineNum), Character: 0},
						End:   protocol.Position{Line: protocol.UInteger(lineNum), Character: protocol.UInteger(len(line) - 1)},
					},
				})
			}

			lineNum++
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return locations, nil
}
