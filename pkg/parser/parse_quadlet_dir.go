package parser

import (
	"encoding/json"
	"os"
	"path"
	"slices"
	"strings"
)

// ParseQuadletDir List all files in the directory and if found any Quadlet
// then parse it.
func ParseQuadletDir(rootDir string) (QuadletDirectory, error) {
	qd := QuadletDirectory{}
	qd.Quadlets = make(map[string]Quadlet)

	// Check for global disable rules
	content, err := os.ReadFile(path.Join(rootDir, ".quadletrc.json"))
	if err == nil {
		disableRules := struct {
			Disabled []string `json:"disabled"`
		}{}
		_ = json.Unmarshal(content, &disableRules)
		qd.DisabledQSR = disableRules.Disabled
	}

	// Check Quadlet files in the directory. Ignore dropins because that is
	// done by the called function.
	quadletExtensions := []string{
		"image", "container", "volume", "network", "kube", "pod", "build",
	}

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return QuadletDirectory{}, err
	}

	for _, e := range entries {
		extension := e.Name()[strings.LastIndex(e.Name(), ".")+1:]

		if !slices.Contains(quadletExtensions, extension) {
			continue
		}

		q, err := ParseQuadlet(ParseQuadletConfig{
			FileName:      e.Name(),
			RootDirectory: rootDir,
		})
		if err != nil {
			return QuadletDirectory{}, err
		}

		qd.Quadlets[e.Name()] = q
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
