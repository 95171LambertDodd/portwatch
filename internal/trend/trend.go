// Package trend tracks port binding frequency over time to detect unusual activity.
package trend

import (
	"sync"
	"time"
)

// Point represents a single observation at a point in time.
type Point struct {
	At    time.Time
	Count int
}

// Tracker records scan counts per key and exposes a recent rate.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	points  map[string][]Point
	clock   func() time.Time
}

// New returns a Tracker that retains data within the given window.
func New(window time.Duration) *Tracker {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock func() time.Time) *Tracker {
	return &Tracker{
		window: window,
		points: make(map[string][]Point),
		clock:  clock,
	}
}

// Record adds a count observation for the given key.
func (t *Tracker) Record(key string, count int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	t.points[key] = append(t.prune(key, now), Point{At: now, Count: count})
}

// Rate returns the average count per window duration for the given key.
// Returns 0 if there are fewer than 2 data points.
func (t *Tracker) Rate(key string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	pts := t.prune(key, now)
	t.points[key] = pts
	if len(pts) < 2 {
		return 0
	}
	var sum int
	for _, p := range pts {
		sum += p.Count
	}
	return float64(sum) / float64(len(pts))
}

// prune removes points outside the retention window (caller must hold mu).
func (t *Tracker) prune(key string, now time.Time) []Point {
	cutoff := now.Add(-t.window)
	pts := t.points[key]
	var kept []Point
	for _, p := range pts {
		if p.At.After(cutoff) {
			kept = append(kept, p)
		}
	}
	return kept
}
