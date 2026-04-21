package aggregator

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/portscanner"
)

var t0 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func makeEntries(ports ...int) []portscanner.PortEntry {
	out := make([]portscanner.PortEntry, len(ports))
	for i, p := range ports {
		out[i] = portscanner.PortEntry{Port: p, Protocol: "tcp", State: "LISTEN"}
	}
	return out
}

func TestNew_ValidParams_ReturnsNonNil(t *testing.T) {
	a, err := New(time.Minute, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestNew_ZeroWindow_ReturnsError(t *testing.T) {
	_, err := New(0, 10)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_ZeroMaxBuckets_ReturnsError(t *testing.T) {
	_, err := New(time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero maxBuckets")
	}
}

func TestAdd_IncrementsCurrentCount(t *testing.T) {
	a, _ := newWithClock(time.Minute, 5, fixedClock(t0))
	a.Add(makeEntries(80, 443))
	if got := a.CurrentCount(); got != 2 {
		t.Errorf("expected count 2, got %d", got)
	}
}

func TestAdd_RotatesBucketAfterWindow(t *testing.T) {
	now := t0
	clock := func() time.Time { return now }
	a, _ := newWithClock(time.Minute, 5, clock)

	a.Add(makeEntries(80))
	now = t0.Add(2 * time.Minute) // advance past window
	a.Add(makeEntries(443))

	buckets := a.Flush()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 completed bucket, got %d", len(buckets))
	}
	if buckets[0].Count != 1 {
		t.Errorf("expected bucket count 1, got %d", buckets[0].Count)
	}
}

func TestFlush_ClearsBuckets(t *testing.T) {
	now := t0
	clock := func() time.Time { return now }
	a, _ := newWithClock(time.Minute, 5, clock)

	a.Add(makeEntries(22))
	now = t0.Add(2 * time.Minute)
	a.Add(makeEntries(80))

	a.Flush()
	if got := a.Flush(); len(got) != 0 {
		t.Errorf("expected empty flush after clear, got %d buckets", len(got))
	}
}

func TestAdd_RespectsMaxBuckets(t *testing.T) {
	now := t0
	clock := func() time.Time { return now }
	a, _ := newWithClock(time.Minute, 2, clock)

	for i := 0; i < 5; i++ {
		a.Add(makeEntries(80 + i))
		now = now.Add(2 * time.Minute)
	}

	buckets := a.Flush()
	if len(buckets) > 2 {
		t.Errorf("expected at most 2 buckets, got %d", len(buckets))
	}
}
