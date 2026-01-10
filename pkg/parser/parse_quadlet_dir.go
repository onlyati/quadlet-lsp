package parser

import (
	"encoding/json"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

// ParseQuadletDir List all files in the directory and if found any Quadlet
// then parse it.
func ParseQuadletDir(rootDir string, dirLevel int) (QuadletDirectory, error) {
	qd := QuadletDirectory{}
	qd.Quadlets = make(map[string]Quadlet)

	// Check for global disable rules
	content, err := os.ReadFile(path.Join(rootDir, ".quadletrc.json"))
	if err == nil {
		disableRules := struct {
			Disabled []string `json:"disable"`
		}{}
		_ = json.Unmarshal(content, &disableRules)
		qd.DisabledQSR = disableRules.Disabled
	}

	// Check Quadlet files in the directory. Ignore dropins because that is
	// done by the called function.
	quadletExtensions := []string{
		"image", "container", "volume", "network", "kube", "pod", "build",
	}

	err = utils.QuadletWalkDir(rootDir, dirLevel, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		extesion := d.Name()[strings.LastIndex(d.Name(), ".")+1:]
		if !slices.Contains(quadletExtensions, extesion) {
			return nil
		}

		fName, _ := strings.CutPrefix(p, rootDir+"/")
		q, err := ParseQuadlet(ParseQuadletConfig{
			FileName:       fName,
			RootDirectory:  rootDir,
			CollectDropins: false,
		})
		if err != nil {
			return err
		}

		qd.Quadlets[fName] = q
		return nil
	})
	if err != nil {
		return qd, err
	}

	// Now check for all dropins
	err = utils.QuadletWalkDir(rootDir, dirLevel, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(d.Name(), ".conf") {
			return nil
		}

		fName, _ := strings.CutPrefix(p, rootDir+"/")
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

		found := false
		for _, sq := range qd.Quadlets {
			tmpPaths := strings.Split(sq.Name, string(os.PathSeparator))
			if isDropinsBelongsToQuadlet(tmpPaths[len(tmpPaths)-1], parentDir) {
				sq.Dropins = append(sq.Dropins, dropin)
				qd.Quadlets[sq.Name] = sq
				found = true
			}
		}
		if !found {
			qd.OprhanDropins = append(qd.OprhanDropins, dropin)
		}

		return nil
	})
	if err != nil {
		return qd, err
	}

	// Now check for the references
	for _, q := range qd.Quadlets {
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
	}

	// Check which quadlet refer to specific one (reverse of References)
	for name := range qd.Quadlets {
		for _, sq := range qd.Quadlets {
			if slices.Contains(sq.References, name) {
				q := qd.Quadlets[name]
				q.PartOf = append(q.PartOf, sq.Name)
				qd.Quadlets[name] = q
			}
		}
	}

	return qd, nil
}
