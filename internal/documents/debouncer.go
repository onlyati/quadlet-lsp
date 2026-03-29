package documents

import (
	"sync"
	"time"
)

// DocumentDebouncer manages timers for individual files
type DocumentDebouncer struct {
	mu     sync.Mutex
	timers map[string]*time.Timer
}

func NewDocumentDebouncer() *DocumentDebouncer {
	return &DocumentDebouncer{
		timers: make(map[string]*time.Timer),
	}
}

// Debounce executes the action after the delay. If called again for the same URI,
// the timer resets.
func (d *DocumentDebouncer) Debounce(uri string, delay time.Duration, action func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Stop existing timer for this specific file
	if timer, exists := d.timers[uri]; exists {
		timer.Stop()
	}

	// Create a new timer
	d.timers[uri] = time.AfterFunc(delay, func() {
		// Clean up the timer from the map when it executes
		d.mu.Lock()
		delete(d.timers, uri)
		d.mu.Unlock()

		action()
	})
}
