// Package dedup provides event deduplication based on a fingerprint cache.
package dedup

import (
	"sync"
	"time"
)

// Entry represents a deduplicated event key with an expiry.
type entry struct {
	seenAt  time.Time
	expires time.Time
}

// Deduplicator suppresses duplicate events within a TTL window.
type Deduplicator struct {
	mu    sync.Mutex
	cache map[string]entry
	ttl   time.Duration
	clock func() time.Time
}

// New returns a Deduplicator with the given TTL.
func New(ttl time.Duration) *Deduplicator {
	return newWithClock(ttl, time.Now)
}

func newWithClock(ttl time.Duration, clock func() time.Time) *Deduplicator {
	return &Deduplicator{
		cache: make(map[string]entry),
		ttl:   ttl,
		clock: clock,
	}
}

// IsDuplicate returns true if the key was seen within the TTL window.
// If not a duplicate, the key is recorded.
func (d *Deduplicator) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.clock()
	if e, ok := d.cache[key]; ok && now.Before(e.expires) {
		return true
	}
	d.cache[key] = entry{seenAt: now, expires: now.Add(d.ttl)}
	return false
}

// Evict removes expired entries from the cache.
func (d *Deduplicator) Evict() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.clock()
	removed := 0
	for k, e := range d.cache {
		if now.After(e.expires) {
			delete(d.cache, k)
			removed++
		}
	}
	return removed
}

// Size returns the number of tracked keys.
func (d *Deduplicator) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.cache)
}
