package utils

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
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
	FilePath   string
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
	DirLevel      int    // How deep we want to search
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
	lines := []QuadletLine{}
	err := QuadletWalkDir(params.RootDirectory, params.DirLevel, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(d.Name(), ".conf") {
			return nil
		}

		// Check what is parent's directory name, it should match with the one on the
		// parameters.
		pathSplit := strings.Split(path, string(os.PathSeparator))
		if len(pathSplit) < 2 {
			return nil
		}
		if pathSplit[len(pathSplit)-2] == dirPath {
			f, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			qs := readItems(string(f), params.Property, params.Section)
			for _, q := range qs {
				q.FilePath = path
				lines = append(lines, q)
			}
		}

		return nil
	})
	if err != nil {
		return nil
	}
	return lines
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

func findImageInContainerUnit(f []byte, props FindImageExposedPortsProperty) []string {
	var images []string

	// If file ends with .conf, this is a dropins.
	// Start the search from the corresponding container
	if strings.HasSuffix(props.URI, ".conf") {
		tmp := strings.Split(props.URI, string(os.PathSeparator))
		parentDirName := tmp[len(tmp)-2]
		ownerContainer, _ := strings.CutSuffix(parentDirName, ".d")
		props.URI = ownerContainer
		newF, err := os.ReadFile(path.Join(props.RootDir, props.URI))
		if err != nil {
			return []string{}
		}
		f = newF
	}
	lines := FindItems(
		FindItemProperty{
			RootDirectory: props.RootDir,
			Text:          string(f),
			Section:       "[Container]",
			Property:      "Image",
			URI:           props.URI,
			DirLevel:      props.DirLevel,
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
					RootDirectory: props.RootDir,
					Text:          string(f),
					Section:       "[Image]",
					Property:      "Image",
					URI:           props.URI,
					DirLevel:      props.DirLevel,
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

type FindImageExposedPortsProperty struct {
	C        Commander
	Name     string
	RootDir  string
	URI      string
	DirLevel int
}

// FindImageExposedPorts This function looking around the current working directory and looking
// for references of the specified name
func FindImageExposedPorts(props FindImageExposedPortsProperty) []string {
	var ports []string

	name := strings.TrimPrefix(props.Name, "file://")
	var images []string

	isItContainerDropins := strings.HasSuffix(name, ".conf") && strings.Contains(name, ".container.d")
	if strings.HasSuffix(name, ".container") || isItContainerDropins {
		f, err := os.ReadFile(name)
		if err != nil {
			log.Printf("failed to read file: %+v", err.Error())
			return ports
		}
		images = findImageInContainerUnit(f, props)
	}

	if strings.HasSuffix(name, ".pod") {
		tmp := strings.Split(name, string(os.PathSeparator))
		podFileName := tmp[len(tmp)-1]
		refs, err := FindReferences(
			GoReferenceProperty{
				Property: "Pod",
				SearchIn: []string{"container"},
				DirLevel: props.DirLevel,
			},
			podFileName,
			props.RootDir,
		)
		for _, r := range refs {
			f, err := os.ReadFile(strings.TrimPrefix(r.URI, "file://"))
			if err != nil {
				continue
			}
			lines := FindItems(
				FindItemProperty{
					RootDirectory: props.RootDir,
					Text:          string(f),
					Section:       "[Container]",
					Property:      "Pod",
					URI:           r.URI,
				},
			)
			if len(lines) > 0 {
				if lines[0].Value == podFileName {
					tmp := findImageInContainerUnit(
						f,
						FindImageExposedPortsProperty{
							C:        props.C,
							Name:     r.URI,
							RootDir:  props.RootDir,
							URI:      r.URI,
							DirLevel: props.DirLevel,
						},
					)
					images = append(images, tmp...)
				}
			}
		}

		if err != nil {
			log.Printf("failed to walk to find container files: %+v", err.Error())
			return ports
		}
	}

	// We have the images, check the exposed ports
	for _, img := range images {
		output, err := props.C.Run(
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

type GoReferenceProperty struct {
	Property string
	SearchIn []string
	DirLevel int
}

// FindReferences in other Quadlet files for the specific file
func FindReferences(prop GoReferenceProperty, currentFileName, rootDir string) ([]protocol.Location, error) {
	var locations []protocol.Location
	files := []protocol.CompletionItem{}

	for _, d := range prop.SearchIn {
		filesTmp, err := ListQuadletFiles(d, rootDir, prop.DirLevel)
		if err != nil {
			return nil, err
		}
		files = append(files, filesTmp...)
	}

	for _, f := range files {
		p := ""
		switch v := f.Documentation.(type) {
		case string:
			p = v
		default:
			return nil, errors.New("unexpected error: documentation is not string")
		}

		p, _ = strings.CutPrefix(p, "From work directory: ")
		content, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}

		fname := p[strings.LastIndex(p, string(os.PathSeparator))+1:]
		section := FirstCharacterToUpper(fname[strings.LastIndex(fname, ".")+1:])

		items := FindItems(FindItemProperty{
			URI:           p,
			RootDirectory: rootDir,
			Text:          string(content),
			Section:       "[" + section + "]",
			Property:      prop.Property,
			DirLevel:      prop.DirLevel,
		})
		for _, item := range items {
			if prop.Property == "Volume" {
				volParts := strings.Split(item.Value, ":")
				item.Value = volParts[0]
			}
			if strings.Contains(item.Value, "@") {
				// If contains '@' then it is a systemd template
				item.Value = ConvertTemplateNameToFile(item.Value)
			}
			if item.Value == currentFileName {
				uri := p
				if item.FilePath != "" {
					uri = item.FilePath
				}
				locations = append(locations, protocol.Location{
					URI: protocol.DocumentUri("file://" + uri),
					Range: protocol.Range{
						Start: protocol.Position{Line: item.LineNumber, Character: 0},
						End:   protocol.Position{Line: item.LineNumber, Character: item.Length},
					},
				})
			}
		}
	}

	return locations, nil
}
