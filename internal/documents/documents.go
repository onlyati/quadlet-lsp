// Package documents reads and parsing continousily the content of the Quadlet files.
package documents

import (
	"strings"
	"sync"

	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
)

type Documents struct {
	mu      sync.RWMutex
	files   map[string]string
	parsers map[string]parser.Parser
}

func NewDocuments() Documents {
	return Documents{
		files:   make(map[string]string),
		parsers: make(map[string]parser.Parser),
	}
}

func (d *Documents) ListFileNames() []string {
	d.mu.RLock()
	keys := make([]string, 0, len(d.files))
	for k := range d.files {
		keys = append(keys, k)
	}
	d.mu.RUnlock()

	return keys
}

func (d *Documents) Add(uri, text string) {
	d.mu.Lock()
	uri = strings.TrimPrefix(uri, "file://")
	d.files[uri] = text
	d.parsers[uri] = parser.NewParser(uri)
	d.mu.Unlock()
}

func (d *Documents) Delete(uri string) {
	d.mu.Lock()
	uri = strings.TrimPrefix(uri, "file://")
	delete(d.files, uri)
	delete(d.parsers, uri)
	d.mu.Unlock()
}

func (d *Documents) Read(uri string) string {
	d.mu.RLock()
	uri = strings.TrimPrefix(uri, "file://")
	text := d.files[uri]
	d.mu.RUnlock()

	return text
}

func (d *Documents) ReadQuadlet(uri string) *parser.QuadletNode {
	d.mu.RLock()
	uri = strings.TrimPrefix(uri, "file://")
	text := d.parsers[uri].Quadlet
	d.mu.RUnlock()

	return text
}

func (d *Documents) CheckURI(uri string) (string, bool) {
	d.mu.RLock()
	uri = strings.TrimPrefix(uri, "file://")
	text, ok := d.files[uri]
	d.mu.RUnlock()

	return text, ok
}
