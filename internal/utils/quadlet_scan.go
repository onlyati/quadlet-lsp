package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type QuadletLine struct {
	LineNumber uint32
	Length     uint32
	Property   string
	Value      string
	RawLine    string
}

type ScanProperty struct {
	Property string
	Section  string
}

func ScanQadlet(
	text string,
	podmanVer PodmanVersion,
	properties map[ScanProperty]struct{},
	action func(q QuadletLine, p PodmanVersion) *protocol.Diagnostic,
) []protocol.Diagnostic {
	var returnValue []protocol.Diagnostic

	currentSection := ""
	sectionRegexp := regexp.MustCompile(`^\[([A-Za-z0-9_.-]+)\]$`)

	for i, rawLine := range strings.Split(text, "\n") {
		line := strings.TrimSpace(rawLine)

		if currentSection != "" && strings.Contains(line, "=") {
			lineArray := strings.SplitN(line, "=", 2)
			if len(lineArray) > 1 {
				_, found := properties[ScanProperty{
					Property: lineArray[0],
					Section:  currentSection,
				}]

				if found {
					diag := action(QuadletLine{
						LineNumber: uint32(i),
						Length:     uint32(len(line)),
						Property:   lineArray[0],
						Value:      lineArray[1],
						RawLine:    rawLine,
					}, podmanVer)
					if diag != nil {
						returnValue = append(returnValue, *diag)
					}
				}
			}
		}

		if sectionRegexp.MatchString(line) {
			currentSection = line
			continue
		}
	}

	return returnValue
}

// This function scanning the passed text and
// looking for property in specific section.
func FindItems(text, section, property string) []QuadletLine {
	var findings []QuadletLine

	inSection := false

	for i, rawLine := range strings.Split(text, "\n") {
		line := strings.TrimSpace(rawLine)

		if inSection && strings.Contains(line, "=") {
			tmp := strings.SplitN(line, "=", 2)
			if len(tmp) > 1 {
				if tmp[0] == property {
					findings = append(findings, QuadletLine{
						LineNumber: uint32(i),
						Length:     uint32(len(line)),
						Property:   tmp[0],
						Value:      tmp[1],
						RawLine:    rawLine,
					})
				}
			}
			continue
		}

		if strings.HasPrefix(line, "[") && line != section {
			inSection = false
			continue
		}

		if line == section {
			inSection = true
			continue
		}
	}

	return findings
}

func findImageInContainerUnit(f []byte) []string {
	var images []string

	lines := FindItems(
		string(f),
		"[Container]",
		"Image",
	)

	for _, line := range lines {
		if strings.HasSuffix(line.Value, ".image") {
			f, err := os.ReadFile(line.Value)
			if err != nil {
				return images
			}
			lines := FindItems(
				string(f),
				"[Image]",
				"Image",
			)

			for _, l := range lines {
				img, _ := strings.CutSuffix(l.Value, ":")
				images = append(images, img)
			}
			continue
		}

		if strings.HasSuffix(line.Value, ".build") {
			// ignore for now
			continue
		}

		// Here a pure image is defined like `Image=something.icr.io/org/name:tag`
		img, _ := strings.CutSuffix(line.Value, ":")

		images = append(images, img)
	}
	return images
}

// This function looking around the current working directory and looking
// for references of the specified name
func FindImageExposedPorts(c Commander, name string) []string {
	var ports []string

	name, _ = strings.CutPrefix(name, "file://")
	var images []string

	if strings.HasSuffix(name, ".container") {
		f, err := os.ReadFile(name)
		if err != nil {
			log.Printf("failed to read file: %+v", err.Error())
			return ports
		}
		images = findImageInContainerUnit(f)
	}

	if strings.HasSuffix(name, ".pod") {
		tmp := strings.Split(name, string(os.PathSeparator))
		podFileName := tmp[len(tmp)-1]

		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !strings.HasSuffix(path, ".container") {
				return nil
			}

			f, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			lines := FindItems(
				string(f),
				"[Container]",
				"Pod",
			)

			if len(lines) > 0 {
				if lines[0].Value == podFileName {
					log.Printf("looking for in %s", lines[0].Value)
					tmp := findImageInContainerUnit(f)
					images = append(images, tmp...)
				}
			}

			return nil
		})
		if err != nil {
			log.Printf("failed to walk to find container files: %+v", err.Error())
			return ports
		}
	}

	// We have the images, check the exposed ports
	for _, img := range images {
		output, err := c.Run(
			"podman",
			"image", "inspect", img,
		)
		if err != nil {
			ports = append(ports, "failed-check-"+img)
			log.Printf("failed to inspect image: %+v", err.Error())
			return ports
		}

		if len(output) == 0 {
			log.Println("image is not pulled")
			continue
		}

		inspectJSON := strings.Join(output, "")
		var data []map[string]any
		json.Unmarshal([]byte(inspectJSON), &data)

		config, ok := data[0]["Config"].(map[string]any)
		if !ok {
			return ports
		}

		exposedPorts, ok := config["ExposedPorts"].(map[string]any)
		if !ok {
			return ports
		}

		for port := range exposedPorts {
			tmp := strings.Split(port, "/")
			ports = append(ports, tmp[0])
		}
	}

	return ports
}

// This function looking for that the cursor currently in which section.
// Sections are like `[Container]`, `[Unit]`, and so on.
func FindSection(lines []string, lineNumber uint32) string {
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
