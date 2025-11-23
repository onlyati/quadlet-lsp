package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
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
	Section    string
}

type ScanProperty struct {
	Property string
	Section  string
}

func ScanQadlet(
	text string,
	podmanVer PodmanVersion,
	properties map[ScanProperty]struct{},
	action func(q QuadletLine, p PodmanVersion, extraInfo any) []protocol.Diagnostic,
	extraInfo any,
) []protocol.Diagnostic {
	var returnValue []protocol.Diagnostic

	currentSection := ""
	sectionRegexp := regexp.MustCompile(`^\[([A-Za-z]+)\]$`)

	// If properts[*] = "*", it means scan all lines
	_, scanAllLines := properties[ScanProperty{Section: "*", Property: "*"}]

	readContinue, lineValue := false, QuadletLine{}

	for i, rawLine := range strings.Split(text, "\n") {
		line := strings.TrimSpace(rawLine)

		// skip if commented line
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if readContinue {
			partialValue, contSign := strings.CutSuffix(rawLine, " \\")
			lineValue.Value += " " + strings.TrimSpace(partialValue)
			if !contSign {
				readContinue = false
				diag := action(lineValue, podmanVer, extraInfo)
				if diag != nil {
					returnValue = append(returnValue, diag...)
				}
			}
			continue
		}

		if currentSection != "" && strings.Contains(line, "=") {
			lineArray := strings.SplitN(line, "=", 2)
			if len(lineArray) > 1 {
				_, found := properties[ScanProperty{
					Property: lineArray[0],
					Section:  currentSection,
				}]

				if found || scanAllLines {
					partialValue, mustContinue := strings.CutSuffix(lineArray[1], " \\")
					lineValue = QuadletLine{
						LineNumber: uint32(i),
						Length:     uint32(len(line)),
						Property:   lineArray[0],
						Value:      strings.TrimSpace(partialValue),
						RawLine:    rawLine,
						Section:    currentSection,
					}

					if mustContinue {
						readContinue = true
						continue
					}

					diag := action(lineValue, podmanVer, extraInfo)
					if diag != nil {
						returnValue = append(returnValue, diag...)
					}
				}
			}
		}

		if sectionRegexp.MatchString(line) {
			currentSection = line
			if scanAllLines {
				diag := action(QuadletLine{
					LineNumber: uint32(i),
					Length:     uint32(len(line)),
					Property:   "",
					Value:      "",
					RawLine:    rawLine,
					Section:    currentSection,
				}, podmanVer, extraInfo)
				if diag != nil {
					returnValue = append(returnValue, diag...)
				}
			}
			continue
		}
	}

	return returnValue
}

// FindItemProperty Parameter of FindItems() function.
type FindItemProperty struct {
	URI           string // URI that is passed from editor
	RootDirectory string // Workspace root directory
	Text          string // Text of the current document
	Section       string // Section we are looking for
	Property      string // Property we are looking for within the section
}

// FindItems This function scanning the passed text and looking for property in specific section.
func FindItems(params FindItemProperty) []QuadletLine {
	fileName := params.URI[strings.LastIndex(params.URI, "/")+1:]
	extension := fileName[strings.LastIndex(fileName, ".")+1:]
	fileName = fileName[:strings.LastIndex(fileName, ".")]

	// First looking for the unit related dropins
	parts := strings.Split(fileName, "-")
	for i := range parts {
		p := parts[0 : len(parts)-i]
		dirPath := strings.Join(p, "-")

		if i != 0 {
			dirPath += "-"
		}
		dirPath += "." + extension + ".d"

		q := findItemsInDir(params, dirPath)
		if len(q) > 0 {
			return q
		}
	}

	// Then looking for the generic <extension>.d directory
	dirPath := path.Join(params.RootDirectory, extension+".d")
	q := findItemsInDir(params, dirPath)
	if len(q) > 0 {
		return q
	}

	// Looking for the in file
	return readItems(params.Text, params.Property, params.Section)
}

func findItemsInDir(params FindItemProperty, dirPath string) []QuadletLine {
	entries, err := os.ReadDir(dirPath)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}

			if !strings.HasSuffix(e.Name(), ".conf") {
				continue
			}

			file, err := os.ReadFile(path.Join(dirPath, e.Name()))
			if err != nil {
				fmt.Printf("failed to open file: %s", err.Error())
				return nil
			}

			q := readItems(string(file), params.Property, params.Section)
			if len(q) > 0 {
				return q
			}
		}
	}

	return nil
}

func readItems(text, property, section string) []QuadletLine {
	var findings []QuadletLine
	inSection := false

	// Value can be split to multiple line using ' \'
	readContinue := false

	for i, rawLine := range strings.Split(text, "\n") {
		line := strings.TrimSpace(rawLine)

		if readContinue {
			addValue, contSign := strings.CutSuffix(rawLine, " \\")
			if !contSign {
				readContinue = false
			}
			findings[len(findings)-1].Value += " " + strings.TrimSpace(addValue)
			continue
		}

		if inSection && strings.Contains(line, "=") {
			tmp := strings.SplitN(line, "=", 2)
			if len(tmp) > 1 {
				if tmp[0] == property {
					value, mustContinue := strings.CutSuffix(tmp[1], " \\")
					readContinue = mustContinue
					findings = append(findings, QuadletLine{
						LineNumber: uint32(i),
						Length:     uint32(len(line)),
						Property:   tmp[0],
						Value:      strings.TrimSpace(value),
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

func findImageInContainerUnit(f []byte, rootDir, uri string) []string {
	var images []string

	lines := FindItems(
		FindItemProperty{
			RootDirectory: rootDir,
			Text:          string(f),
			Section:       "[Container]",
			Property:      "Image",
			URI:           uri,
		},
	)

	for _, line := range lines {
		if strings.HasSuffix(line.Value, ".image") {
			f, err := os.ReadFile(line.Value)
			if err != nil {
				return images
			}
			lines := FindItems(
				FindItemProperty{
					RootDirectory: rootDir,
					Text:          string(f),
					Section:       "[Image]",
					Property:      "Image",
					URI:           uri,
				},
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

// FindImageExposedPorts This function looking around the current working directory and looking
// for references of the specified name
func FindImageExposedPorts(c Commander, name, rootDir, uri string) []string {
	var ports []string

	name, _ = strings.CutPrefix(name, "file://")
	var images []string

	if strings.HasSuffix(name, ".container") {
		f, err := os.ReadFile(name)
		if err != nil {
			log.Printf("failed to read file: %+v", err.Error())
			return ports
		}
		images = findImageInContainerUnit(f, rootDir, uri)
	}

	if strings.HasSuffix(name, ".pod") {
		tmp := strings.Split(name, string(os.PathSeparator))
		podFileName := tmp[len(tmp)-1]

		err := filepath.Walk(".", func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			isItQuadletFile := strings.HasSuffix(p, ".container")
			if info.IsDir() || !isItQuadletFile {
				return nil
			}

			f, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			lines := FindItems(
				FindItemProperty{
					RootDirectory: rootDir,
					Text:          string(f),
					Section:       "[Container]",
					Property:      "Pod",
					URI:           "file://" + path.Join(rootDir, info.Name()),
				},
			)

			if len(lines) > 0 {
				if lines[0].Value == podFileName {
					tmp := findImageInContainerUnit(
						f,
						rootDir,
						"file://"+path.Join(rootDir, info.Name()),
					)
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

// FindSection This function looking for that the cursor currently in which section.
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
