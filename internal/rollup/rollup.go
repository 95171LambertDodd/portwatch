// Package rollup groups repeated alerts into summarized events
// to reduce noise when the same port binding fires many times.
package rollup

import (
	"fmt"
	"sync"
	"time"
)

// Event holds a summarized rollup of repeated occurrences.
type Event struct {
	Key       string
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Aggregator accumulates events and flushes summaries after a window.
type Aggregator struct {
	mu     sync.Mutex
	window time.Duration
	clock  func() time.Time
	bucket map[string]*Event
}

// New returns an Aggregator that groups events within the given window.
func New(window time.Duration) *Aggregator {
	return &Aggregator{
		window: window,
		clock:  time.Now,
		bucket: make(map[string]*Event),
	}
}

func newWithClock(window time.Duration, clock func() time.Time) *Aggregator {
	a := New(window)
	a.clock = clock
	return a
}

// Record adds an occurrence for the given key.
func (a *Aggregator) Record(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	now := a.clock()
	if e, ok := a.bucket[key]; ok {
		e.Count++
		e.LastSeen = now
	} else {
		a.bucket[key] = &Event{Key: key, Count: 1, FirstSeen: now, LastSeen: now}
	}
}

// Flush returns all events whose window has elapsed and removes them.
func (a *Aggregator) Flush() []Event {
	a.mu.Lock()
	defer a.mu.Unlock()
	now := a.clock()
	var out []Event
	for k, e := range a.bucket {
		if now.Sub(e.FirstSeen) >= a.window {
			out = append(out, *e)
			delete(a.bucket, k)
		}
	}
	return out
}

// Summary returns a human-readable string for an event.
func Summary(e Event) string {
	return fmt.Sprintf("[rollup] key=%s count=%d first=%s last=%s",
		e.Key, e.Count, e.FirstSeen.Format(time.RFC3339), e.LastSeen.Format(time.RFC3339))
}
