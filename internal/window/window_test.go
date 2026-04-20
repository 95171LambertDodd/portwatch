package window

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestNew_PositiveWindow_ReturnsNonNil(t *testing.T) {
	c, err := New(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestNew_ZeroDuration_ReturnsError(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNew_NegativeDuration_ReturnsError(t *testing.T) {
	_, err := New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestRecord_And_Count_WithinWindow(t *testing.T) {
	now := time.Unix(1000, 0)
	c, _ := newWithClock(10*time.Second, fixedClock(now))

	c.Record()
	c.Record()
	c.Record()

	if got := c.Count(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestCount_EvictsExpiredEvents(t *testing.T) {
	base := time.Unix(1000, 0)
	var current time.Time
	clock := func() time.Time { return current }

	c, _ := newWithClock(5*time.Second, clock)

	// Record 3 events at t=0
	current = base
	c.Record()
	c.Record()
	c.Record()

	// Advance past the window
	current = base.Add(6 * time.Second)
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", got)
	}
}

func TestCount_PartialEviction(t *testing.T) {
	base := time.Unix(1000, 0)
	var current time.Time
	clock := func() time.Time { return current }

	c, _ := newWithClock(10*time.Second, clock)

	current = base
	c.Record() // will expire

	current = base.Add(8 * time.Second)
	c.Record() // still within window later

	current = base.Add(11 * time.Second)
	if got := c.Count(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestReset_ClearsAllEvents(t *testing.T) {
	now := time.Unix(1000, 0)
	c, _ := newWithClock(10*time.Second, fixedClock(now))

	c.Record()
	c.Record()
	c.Reset()

	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
