// Package quota enforces per-key event rate quotas over a sliding window.
// It tracks how many events a given key has produced and rejects further
// events once the configured limit is reached within the window duration.
package quota

import (
	"errors"
	"sync"
	"time"
)

// Clock allows injecting a fake time source in tests.
type Clock func() time.Time

// Enforcer tracks event counts per key and enforces a maximum quota.
type Enforcer struct {
	mu      sync.Mutex
	window  time.Duration
	limit   int
	clock   Clock
	buckets map[string]*bucket
}

type bucket struct {
	count     int
	windowEnd time.Time
}

// New creates an Enforcer with the given sliding window and per-key limit.
// Returns an error if window is zero/negative or limit is non-positive.
func New(window time.Duration, limit int) (*Enforcer, error) {
	return newWithClock(window, limit, time.Now)
}

func newWithClock(window time.Duration, limit int, clock Clock) (*Enforcer, error) {
	if window <= 0 {
		return nil, errors.New("quota: window must be positive")
	}
	if limit <= 0 {
		return nil, errors.New("quota: limit must be positive")
	}
	return &Enforcer{
		window:  window,
		limit:   limit,
		clock:   clock,
		buckets: make(map[string]*bucket),
	}, nil
}

// Allow returns true and increments the counter if the key is within quota.
// Returns false without incrementing if the quota has been exhausted for the
// current window. The window resets automatically after it expires.
func (e *Enforcer) Allow(key string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.clock()
	b, ok := e.buckets[key]
	if !ok || now.After(b.windowEnd) {
		e.buckets[key] = &bucket{
			count:     1,
			windowEnd: now.Add(e.window),
		}
		return true
	}
	if b.count >= e.limit {
		return false
	}
	b.count++
	return true
}

// Remaining returns how many more events the key may produce in the current
// window. Returns the full limit if no window is active for the key.
func (e *Enforcer) Remaining(key string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.clock()
	b, ok := e.buckets[key]
	if !ok || now.After(b.windowEnd) {
		return e.limit
	}
	rem := e.limit - b.count
	if rem < 0 {
		return 0
	}
	return rem
}
