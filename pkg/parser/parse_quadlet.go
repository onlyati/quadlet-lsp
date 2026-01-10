package parser

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

// ParseQuadletConfig Configuration for ParseQuadlet function
type ParseQuadletConfig struct {
	FileName       string // Name of the file (with extenstion) what should be parsed
	RootDirectory  string // Directory where the Quadlet files are located
	CollectDropins bool   // Want to collect dropins information too?
	DirLevel       int    // How deep in the tree dropins shoold searched
}

// ParseQuadlet This function parse Quadlet file, including its dropins
func ParseQuadlet(c ParseQuadletConfig) (Quadlet, error) {
	q := Quadlet{}
	q.Properties = make(map[string][]QuadletProperty)
	q.Name = c.FileName

	// First parse the file itself
	qPath := fmt.Sprintf("%s%c%s", c.RootDirectory, os.PathSeparator, c.FileName)
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
	// Now check for all dropins

	if c.CollectDropins {
		err = utils.QuadletWalkDir(c.RootDirectory, c.DirLevel, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !strings.HasSuffix(d.Name(), ".conf") {
				return nil
			}

			fName, _ := strings.CutPrefix(p, c.RootDirectory+"/")
			pathsTmp := strings.Split(fName, string(os.PathSeparator))
			if len(pathsTmp) < 2 {
				return nil
			}
			parentDir := pathsTmp[len(pathsTmp)-2]

			f, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			lastSection := ""
			lastLine := ""
			dropin := Dropin{
				FileName:   fName,
				Directory:  parentDir,
				Properties: make(map[string][]QuadletProperty),
				SourceFile: string(f),
			}

			for l := range strings.SplitSeq(string(f), "\n") {
				err := parseQuadletContent(
					l,
					ParseQuadletConfig{},
					&dropin.Properties,
					&lastSection,
					&lastLine,
				)
				if err != nil {
					return err
				}
			}

			if isDropinsBelongsToQuadlet(q.Name, parentDir) {
				q.Dropins = append(q.Dropins, dropin)
			}

			return nil
		})
		if err != nil {
			return q, err
		}
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

	//
	// Convert header markdown to HTML
	//
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(strings.Join(q.Header, "\n")))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	q.HeaderHTML = string(markdown.Render(doc, renderer))

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

		content, err := os.ReadFile(fmt.Sprintf("%s%c%s", dirPath, os.PathSeparator, e.Name()))
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
