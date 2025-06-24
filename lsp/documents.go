package lsp

import "sync"

type Documents struct {
	mu    sync.RWMutex
	files map[string]string
}

func newDocuments() Documents {
	return Documents{
		files: make(map[string]string),
	}
}

func (d *Documents) add(uri, text string) {
	d.mu.Lock()
	d.files[uri] = text
	d.mu.Unlock()
}

func (d *Documents) delete(uri string) {
	d.mu.Lock()
	delete(d.files, uri)
	d.mu.Unlock()
}

func (d *Documents) read(uri string) string {
	d.mu.RLock()
	text := d.files[uri]
	d.mu.RUnlock()

	return text
}

func (d *Documents) checkUri(uri string) (string, bool) {
	d.mu.RLock()
	text, ok := d.files[uri]
	d.mu.RUnlock()

	return text, ok
}
