package utils

import "sync"

type Documents struct {
	mu    sync.RWMutex
	files map[string]string
}

func NewDocuments() Documents {
	return Documents{
		files: make(map[string]string),
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
	d.files[uri] = text
	d.mu.Unlock()
}

func (d *Documents) Delete(uri string) {
	d.mu.Lock()
	delete(d.files, uri)
	d.mu.Unlock()
}

func (d *Documents) Read(uri string) string {
	d.mu.RLock()
	text := d.files[uri]
	d.mu.RUnlock()

	return text
}

func (d *Documents) CheckUri(uri string) (string, bool) {
	d.mu.RLock()
	text, ok := d.files[uri]
	d.mu.RUnlock()

	return text, ok
}
