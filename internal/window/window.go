// Package window provides a sliding time-window counter for tracking
// event frequency over a rolling interval.
package window

import (
	"fmt"
	"sync"
	"time"
)

// Clock allows injecting a time source for testing.
type Clock func() time.Time

// Counter tracks the number of events recorded within a sliding window.
type Counter struct {
	mu         sync.Mutex
	window     time.Duration
	timestamps []time.Time
	clock      Clock
}

// New creates a Counter with the given sliding window duration.
func New(window time.Duration) (*Counter, error) {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock Clock) (*Counter, error) {
	if window <= 0 {
		return nil, fmt.Errorf("window duration must be positive, got %s", window)
	}
	return &Counter{
		window: window,
		clock:  clock,
	}, nil
}

// Record adds an event at the current time.
func (c *Counter) Record() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	c.timestamps = append(c.timestamps, now)
	c.evict(now)
}

// Count returns the number of events recorded within the current window.
func (c *Counter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(c.clock())
	return len(c.timestamps)
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timestamps = c.timestamps[:0]
}

// evict removes timestamps older than the window. Must be called with mu held.
func (c *Counter) evict(now time.Time) {
	cutoff := now.Add(-c.window)
	i := 0
	for i < len(c.timestamps) && c.timestamps[i].Before(cutoff) {
		i++
	}
	c.timestamps = c.timestamps[i:]
}
