// Package debounce provides a debouncer that suppresses repeated events
// within a configurable quiet period, emitting only after the period elapses.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays firing a callback until no new events arrive within the
// quiet window. Each unique key is debounced independently.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	timers  map[string]*time.Timer
	callback func(key string)
}

// New creates a Debouncer with the given quiet window and callback.
// The callback is invoked with the key after no new Trigger calls
// arrive for that key within the window.
func New(window time.Duration, callback func(key string)) *Debouncer {
	return &Debouncer{
		window:   window,
		timers:   make(map[string]*time.Timer),
		callback: callback,
	}
}

// Trigger resets the debounce timer for key. If no further Trigger calls
// arrive within the window, the callback is invoked.
func (d *Debouncer) Trigger(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.window, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		d.callback(key)
	})
}

// Cancel stops any pending timer for key without invoking the callback.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys currently awaiting debounce.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
