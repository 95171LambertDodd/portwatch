package dedup

import (
	"testing"
	"time"
)

var (
	baseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_ReturnsNonNil(t *testing.T) {
	d := New(time.Minute)
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestIsDuplicate_FirstCallNotDuplicate(t *testing.T) {
	d := newWithClock(time.Minute, fixedClock(baseTime))
	if d.IsDuplicate("key1") {
		t.Error("first call should not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinTTL(t *testing.T) {
	d := newWithClock(time.Minute, fixedClock(baseTime))
	d.IsDuplicate("key1")
	if !d.IsDuplicate("key1") {
		t.Error("second call within TTL should be a duplicate")
	}
}

func TestIsDuplicate_CallAfterTTLExpiry(t *testing.T) {
	now := baseTime
	clock := func() time.Time { return now }
	d := newWithClock(time.Minute, clock)
	d.IsDuplicate("key1")
	now = baseTime.Add(2 * time.Minute)
	if d.IsDuplicate("key1") {
		t.Error("call after TTL should not be a duplicate")
	}
}

func TestIsDuplicate_DifferentKeysAreIndependent(t *testing.T) {
	d := newWithClock(time.Minute, fixedClock(baseTime))
	d.IsDuplicate("key1")
	if d.IsDuplicate("key2") {
		t.Error("different keys should be independent")
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	now := baseTime
	clock := func() time.Time { return now }
	d := newWithClock(time.Minute, clock)
	d.IsDuplicate("key1")
	d.IsDuplicate("key2")
	now = baseTime.Add(2 * time.Minute)
	d.IsDuplicate("key3") // fresh
	removed := d.Evict()
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if d.Size() != 1 {
		t.Errorf("expected size 1 after eviction, got %d", d.Size())
	}
}

func TestSize_ReflectsCache(t *testing.T) {
	d := newWithClock(time.Minute, fixedClock(baseTime))
	if d.Size() != 0 {
		t.Error("expected empty cache")
	}
	d.IsDuplicate("a")
	d.IsDuplicate("b")
	if d.Size() != 2 {
		t.Errorf("expected size 2, got %d", d.Size())
	}
}
