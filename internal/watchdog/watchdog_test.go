package watchdog

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_ReturnsNonNil(t *testing.T) {
	w := New(5 * time.Second)
	if w == nil {
		t.Fatal("expected non-nil watchdog")
	}
}

func TestCheck_NoHeartbeat_ReturnsStale(t *testing.T) {
	w := New(5 * time.Second)
	status, msg := w.Check()
	if status != StatusStale {
		t.Errorf("expected stale, got %s", status)
	}
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestCheck_RecentBeat_ReturnsOK(t *testing.T) {
	now := time.Now()
	w := newWithClock(5*time.Second, fixedClock(now))
	w.Beat()
	status, _ := w.Check()
	if status != StatusOK {
		t.Errorf("expected ok, got %s", status)
	}
}

func TestCheck_StaleBeat_ReturnsStale(t *testing.T) {
	base := time.Now()
	current := base
	w := newWithClock(5*time.Second, func() time.Time { return current })
	w.Beat()
	current = base.Add(10 * time.Second)
	status, msg := w.Check()
	if status != StatusStale {
		t.Errorf("expected stale, got %s", status)
	}
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestLastBeat_ReflectsMostRecentBeat(t *testing.T) {
	now := time.Now()
	w := newWithClock(5*time.Second, fixedClock(now))
	if !w.LastBeat().IsZero() {
		t.Error("expected zero before any beat")
	}
	w.Beat()
	if !w.LastBeat().Equal(now) {
		t.Errorf("expected %v, got %v", now, w.LastBeat())
	}
}

func TestBeat_UpdatesOnEachCall(t *testing.T) {
	base := time.Now()
	current := base
	w := newWithClock(5*time.Second, func() time.Time { return current })
	w.Beat()
	first := w.LastBeat()
	current = base.Add(2 * time.Second)
	w.Beat()
	second := w.LastBeat()
	if !second.After(first) {
		t.Errorf("expected second beat after first: %v vs %v", second, first)
	}
}
