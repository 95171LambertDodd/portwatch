// Package eviction provides an LRU-based eviction policy for bounded caches
// used throughout portwatch to prevent unbounded memory growth.
package eviction

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Entry holds a cached value with its expiry metadata.
type Entry struct {
	Key       string
	Value     any
	ExpiresAt time.Time
}

// Cache is a thread-safe LRU cache with optional TTL-based eviction.
type Cache struct {
	mu       sync.Mutex
	cap      int
	ttl      time.Duration
	items    map[string]*list.Element
	order    *list.List
	now      func() time.Time
}

// New creates a Cache with the given capacity and TTL.
// A zero TTL disables time-based expiry.
func New(capacity int, ttl time.Duration) (*Cache, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("eviction: capacity must be positive, got %d", capacity)
	}
	return newWithClock(capacity, ttl, time.Now), nil
}

func newWithClock(capacity int, ttl time.Duration, now func() time.Time) *Cache {
	return &Cache{
		cap:   capacity,
		ttl:   ttl,
		items: make(map[string]*list.Element, capacity),
		order: list.New(),
		now:   now,
	}
}

// Set inserts or updates a key in the cache.
func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiry := time.Time{}
	if c.ttl > 0 {
		expiry = c.now().Add(c.ttl)
	}

	if el, ok := c.items[key]; ok {
		c.order.MoveToFront(el)
		el.Value.(*Entry).Value = value
		el.Value.(*Entry).ExpiresAt = expiry
		return
	}

	if c.order.Len() >= c.cap {
		c.evictOldest()
	}

	entry := &Entry{Key: key, Value: value, ExpiresAt: expiry}
	el := c.order.PushFront(entry)
	c.items[key] = el
}

// Get retrieves a value from the cache. Returns (value, true) on hit.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.items[key]
	if !ok {
		return nil, false
	}

	entry := el.Value.(*Entry)
	if c.ttl > 0 && c.now().After(entry.ExpiresAt) {
		c.remove(el)
		return nil, false
	}

	c.order.MoveToFront(el)
	return entry.Value, true
}

// Len returns the number of items currently in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

func (c *Cache) evictOldest() {
	back := c.order.Back()
	if back != nil {
		c.remove(back)
	}
}

func (c *Cache) remove(el *list.Element) {
	c.order.Remove(el)
	delete(c.items, el.Value.(*Entry).Key)
}
