package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses duplicate alerts for the same key within a cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New creates a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// If allowed, it records the current time for the key.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if t, ok := l.last[key]; ok {
		if now.Sub(t) < l.cooldown {
			return false
		}
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded time for a specific key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Flush removes all recorded entries, resetting the limiter entirely.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
