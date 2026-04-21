// Package aggregator groups port scan events into time-bucketed summaries
// for downstream reporting and trend analysis.
package aggregator

import (
	"fmt"
	"sync"
	"time"

	"github.com/yourorg/portwatch/internal/portscanner"
)

// Bucket holds aggregated port entries for a specific time window.
type Bucket struct {
	Start   time.Time
	End     time.Time
	Entries []portscanner.PortEntry
	Count   int
}

// Aggregator accumulates scan results into fixed-duration buckets.
type Aggregator struct {
	mu         sync.Mutex
	window     time.Duration
	buckets    []Bucket
	current    *Bucket
	maxBuckets int
	clock      func() time.Time
}

// New creates an Aggregator with the given bucket window and max retained buckets.
func New(window time.Duration, maxBuckets int) (*Aggregator, error) {
	return newWithClock(window, maxBuckets, time.Now)
}

func newWithClock(window time.Duration, maxBuckets int, clock func() time.Time) (*Aggregator, error) {
	if window <= 0 {
		return nil, fmt.Errorf("aggregator: window must be positive, got %s", window)
	}
	if maxBuckets <= 0 {
		return nil, fmt.Errorf("aggregator: maxBuckets must be positive, got %d", maxBuckets)
	}
	return &Aggregator{
		window:     window,
		maxBuckets: maxBuckets,
		clock:      clock,
	}, nil
}

// Add inserts port entries into the current time bucket, rotating if needed.
func (a *Aggregator) Add(entries []portscanner.PortEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := a.clock()
	if a.current == nil || now.After(a.current.End) {
		a.rotateLocked(now)
	}
	a.current.Entries = append(a.current.Entries, entries...)
	a.current.Count += len(entries)
}

// Flush returns all completed buckets and clears them from memory.
func (a *Aggregator) Flush() []Bucket {
	a.mu.Lock()
	defer a.mu.Unlock()

	out := make([]Bucket, len(a.buckets))
	copy(out, a.buckets)
	a.buckets = nil
	return out
}

// CurrentCount returns the number of entries in the active bucket.
func (a *Aggregator) CurrentCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.current == nil {
		return 0
	}
	return a.current.Count
}

func (a *Aggregator) rotateLocked(now time.Time) {
	if a.current != nil {
		a.buckets = append(a.buckets, *a.current)
		if len(a.buckets) > a.maxBuckets {
			a.buckets = a.buckets[len(a.buckets)-a.maxBuckets:]
		}
	}
	start := now.Truncate(a.window)
	a.current = &Bucket{
		Start: start,
		End:   start.Add(a.window),
	}
}
