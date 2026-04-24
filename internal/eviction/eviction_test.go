package eviction

import (
	"testing"
	"time"
)

var (
	fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
)

func fixedClock(offset time.Duration) func() time.Time {
	current := fixedNow
	return func() time.Time {
		t := current
		current = current.Add(offset)
		return t
	}
}

func TestNew_ValidCapacity_ReturnsNonNil(t *testing.T) {
	c, err := New(10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestNew_ZeroCapacity_ReturnsError(t *testing.T) {
	_, err := New(0, 0)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestNew_NegativeCapacity_ReturnsError(t *testing.T) {
	_, err := New(-1, 0)
	if err == nil {
		t.Fatal("expected error for negative capacity")
	}
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	c := newWithClock(5, 0, func() time.Time { return fixedNow })
	c.Set("port:8080", "nginx")
	v, ok := c.Get("port:8080")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v != "nginx" {
		t.Fatalf("expected 'nginx', got %v", v)
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	c := newWithClock(5, 0, func() time.Time { return fixedNow })
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestEviction_OldestRemovedWhenFull(t *testing.T) {
	c := newWithClock(3, 0, func() time.Time { return fixedNow })
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("d", 4) // should evict "a"

	_, ok := c.Get("a")
	if ok {
		t.Fatal("expected 'a' to be evicted")
	}
	if c.Len() != 3 {
		t.Fatalf("expected len 3, got %d", c.Len())
	}
}

func TestTTL_ExpiredEntry_ReturnsMiss(t *testing.T) {
	// clock advances 2s per call: Set uses t=0, Get uses t=2s > ttl=1s
	c := newWithClock(5, 1*time.Second, fixedClock(2*time.Second))
	c.Set("key", "value")
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestTTL_ValidEntry_ReturnsHit(t *testing.T) {
	// clock does not advance — TTL never expires
	c := newWithClock(5, 10*time.Second, func() time.Time { return fixedNow })
	c.Set("key", "value")
	v, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit within TTL")
	}
	if v != "value" {
		t.Fatalf("expected 'value', got %v", v)
	}
}

func TestSet_UpdateExistingKey_DoesNotGrow(t *testing.T) {
	c := newWithClock(5, 0, func() time.Time { return fixedNow })
	c.Set("k", "v1")
	c.Set("k", "v2")
	if c.Len() != 1 {
		t.Fatalf("expected len 1 after update, got %d", c.Len())
	}
	v, _ := c.Get("k")
	if v != "v2" {
		t.Fatalf("expected updated value 'v2', got %v", v)
	}
}
