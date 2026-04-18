package throttle

import (
	"sync"
	"time"
)

// Clock allows injecting time in tests.
type Clock func() time.Time

// Throttle limits how often an action can be triggered per key,
// supporting a max burst before enforcing a cooldown window.
type Throttle struct {
	mu       sync.Mutex
	clock    Clock
	window   time.Duration
	maxBurst int
	state    map[string]*bucket
}

type bucket struct {
	count    int
	windowAt time.Time
}

// New creates a Throttle with the given window duration and burst limit.
func New(window time.Duration, maxBurst int, clock Clock) *Throttle {
	if clock == nil {
		clock = time.Now
	}
	return &Throttle{
		clock:    clock,
		window:   window,
		maxBurst: maxBurst,
		state:    make(map[string]*bucket),
	}
}

// Allow returns true if the action for key is permitted under the burst policy.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	b, ok := t.state[key]
	if !ok || now.After(b.windowAt.Add(t.window)) {
		t.state[key] = &bucket{count: 1, windowAt: now}
		return true
	}
	if b.count < t.maxBurst {
		b.count++
		return true
	}
	return false
}

// Reset clears throttle state for a key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.state, key)
}

// Stats returns the current count and window start for a key, or zero values.
func (t *Throttle) Stats(key string) (count int, since time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if b, ok := t.state[key]; ok {
		return b.count, b.windowAt
	}
	return 0, time.Time{}
}
