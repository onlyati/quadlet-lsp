package parser

import (
	"errors"
	"os"
	"path"
	"slices"
	"strings"
)

// ParseQuadletConfig Configuration for ParseQuadlet function
type ParseQuadletConfig struct {
	FileName      string // Name of the file (with extenstion) what should be parsed
	RootDirectory string // Directory where the Quadlet files are located
}

// ParseQuadlet This function parse Quadlet file, including its dropins
func ParseQuadlet(c ParseQuadletConfig) (Quadlet, error) {
	q := Quadlet{}
	q.Properties = make(map[string][]QuadletProperty)
	q.Name = c.FileName

	// First parse the file itself
	qPath := path.Join(c.RootDirectory, c.FileName)
	content, err := os.ReadFile(qPath)
	if err != nil {
		return Quadlet{}, err
	}
	q.SourceFile = string(content)

	recordQSR := true
	recordHeader := false
	lastLine := ""
	lastSection := ""
	for l := range strings.SplitSeq(string(content), "\n") {
		//
		// First parse the header, all comment lines before the first uncomment line
		//
		if !isItComment(l) {
			recordHeader = false
			recordQSR = false
		}

		if recordQSR {
			tl := removeCommentSign(l)
			if v, found := strings.CutPrefix(tl, "disable-qsr:"); found {
				for r := range strings.SplitSeq(v, " ") {
					if r == "" {
						continue
					}
					if !slices.Contains(q.DisabledQSR, r) {
						q.DisabledQSR = append(q.DisabledQSR, r)
					}
				}
			} else {
				recordQSR = false
				recordHeader = true
			}
		}

		if recordHeader {
			tl := removeCommentSign(l)
			q.Header = append(q.Header, tl)
		}

		if recordHeader || recordQSR {
			continue
		}
		//
		// Now parse the effective part of the quadlet
		//
		err := parseQuadletContent(
			l,
			c,
			&q.Properties,
			&lastSection,
			&lastLine)
		if err != nil {
			return Quadlet{}, err
		}
	}

	//
	// Now parse the dropins
	//
	fileName := q.Name[strings.LastIndex(q.Name, "/")+1:]
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
		parentDir := dirPath
		dirPath = path.Join(c.RootDirectory, dirPath)

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			continue // Directory probably just does not exists, skip
		}

		d, err := parseDropins(dirPath, parentDir, entries)
		if err != nil {
			return Quadlet{}, err
		}
		q.Dropins = append(q.Dropins, d...)
	}

	// Then looking for the generic <extension>.d directory
	dirPath := path.Join(c.RootDirectory, extension+".d")
	parentDir := extension + ".d"
	entries, err := os.ReadDir(dirPath)
	if err == nil {
		d, err := parseDropins(dirPath, parentDir, entries)
		if err != nil {
			return Quadlet{}, err
		}
		q.Dropins = append(q.Dropins, d...)
	}

	//
	// Now check for the references
	//
	imageRef := ""
	podRef := ""
	for _, d := range q.Dropins {
		countReferences(
			&q,
			d.Properties,
			&imageRef,
			&podRef,
		)
	}

	countReferences(
		&q,
		q.Properties,
		&imageRef,
		&podRef,
	)
	if imageRef != "" && imageRef != "nolink" {
		q.References = append(q.References, imageRef)
	}
	if podRef != "" {
		q.References = append(q.References, podRef)
	}

	return q, nil
}

// parseQuadletContent This function parse the get lines from Quadlet file
// and dropins config file as well and convert them to a map.
func parseQuadletContent(
	l string,
	c ParseQuadletConfig,
	qs *map[string][]QuadletProperty,
	lastSection *string,
	lastLine *string,
) error {
	l = strings.TrimSpace(l)
	if isItComment(l) || l == "" {
		return nil
	}

	// Check for section header, like '[Container]'
	if section, found := gatherSectionName(l); found {
		*lastSection = section
		*lastLine = ""
		return nil
	}

	// Only parse further if we are in any section
	if *lastSection == "" {
		return nil
	}

	// Parse the propeterties
	parts := strings.SplitN(l, "=", 2)

	if len(parts) != 2 && !strings.HasSuffix(*lastLine, " \\") {
		// Something is wrong, probably incorrect line
		return errors.New(c.FileName + ": invalid line: " + l)
	}

	if strings.HasSuffix(*lastLine, " \\") {
		// Append value to the lasts added value
		lastValue := (*qs)[*lastSection][len((*qs)[*lastSection])-1].Value
		lastValue = strings.TrimSuffix(lastValue, " \\")
		lastValue = strings.TrimSpace(lastValue)
		(*qs)[*lastSection][len((*qs)[*lastSection])-1].Value = lastValue + " " + l
	} else {
		// Add new item
		(*qs)[*lastSection] = append((*qs)[*lastSection], QuadletProperty{
			Property: strings.TrimSpace(parts[0]),
			Value:    strings.TrimSpace(parts[1]),
		})
	}

	*lastLine = l

	return nil
}

// countReferences Calculate references in the file. It does not check that
// file itself exists, just record what kind of reference it has.
// For example 'Volume=foo.volume' or 'Image=foo.image'
func countReferences(
	q *Quadlet,
	entries map[string][]QuadletProperty,
	imageRef *string,
	podRef *string,
) {
	for k, v := range entries {
		if k != "Container" && k != "Pod" && k != "Kube" {
			continue
		}

		for _, vv := range v {
			if vv.Property == "Image" {
				isItBuild := strings.HasSuffix(vv.Value, ".build")
				isItImage := strings.HasSuffix(vv.Value, ".image")
				isAlreadyRead := *imageRef != ""
				if !isAlreadyRead && (isItBuild || isItImage) {
					*imageRef = vv.Value
				} else {
					if *imageRef == "" {
						*imageRef = "nolink"
					}
				}
			}

			if vv.Property == "Pod" {
				isAlreadyRead := *podRef != ""
				if !isAlreadyRead {
					*podRef = vv.Value
				}
			}

			if vv.Property == "Network" {
				isItNetwork := strings.HasSuffix(vv.Value, ".network")
				if isItNetwork && !slices.Contains(q.References, vv.Value) {
					q.References = append(q.References, vv.Value)
				}
			}

			if vv.Property == "Volume" {
				volume := strings.Split(vv.Value, ":")[0]
				isItVolume := strings.HasSuffix(volume, ".volume")
				if isItVolume && !slices.Contains(q.References, volume) {
					q.References = append(q.References, volume)
				}
			}
		}
	}
}

// parseDropins List '*.conf' files in dropins directory then read them and
// parse using parseQuadletContent function.
func parseDropins(dirPath, parentDir string, entries []os.DirEntry) ([]Dropin, error) {
	dropins := []Dropin{}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".conf") {
			continue
		}

		content, err := os.ReadFile(path.Join(dirPath, e.Name()))
		if err != nil {
			return nil, err
		}
		config := ParseQuadletConfig{
			FileName:      e.Name(),
			RootDirectory: dirPath,
		}
		lastSection := ""
		lastLine := ""

		dropin := Dropin{
			FileName:   e.Name(),
			Directory:  parentDir,
			Properties: make(map[string][]QuadletProperty),
			SourceFile: string(content),
		}

		for l := range strings.SplitSeq(string(content), "\n") {
			err := parseQuadletContent(
				l,
				config,
				&dropin.Properties,
				&lastSection,
				&lastLine,
			)
			if err != nil {
				return nil, err
			}
		}
		dropins = append(dropins, dropin)
	}

	return dropins, nil
}
