package lsp

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// This function looking for that the cursor currently in which section.
// Sections are like `[Container]`, `[Unit]`, and so on.
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

// Walk through on each file in a directory and looking for lines
// that starts with the specified parameter. Used, for example, when
// try to find references for a quadlet file.
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

// Recognize that what image is used in the specific quadlet.
func findImageName(lines []string, lineNumber protocol.UInteger) string {
	// First looking for `Image=value` value
	// First looing for reverse, people usually define image first then parameters
	imageName := ""
	for i := lineNumber; i > 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "Image=") {
			tmp := strings.Split(line, "=")
			if len(tmp) != 2 {
				break
			}
			imageName = tmp[1]
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// We've reached the start of section, try in other direction
			break
		}
	}

	// Check rest of the file for `Image=`
	if imageName == "" {
		for i := lineNumber; int(i) < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if strings.HasPrefix(line, "Image=") {
				tmp := strings.Split(line, "=")
				if len(tmp) != 2 {
					break
				}
				imageName = tmp[1]
			}
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				// We've reached the start of another section
				break
			}
		}
	}
	return imageName
}
