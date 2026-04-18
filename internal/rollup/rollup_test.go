package rollup

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_CreatesEvent(t *testing.T) {
	a := newWithClock(time.Minute, fixedClock(epoch))
	a.Record("tcp:8080")
	a.mu.Lock()
	e := a.bucket["tcp:8080"]
	a.mu.Unlock()
	if e == nil || e.Count != 1 {
		t.Fatalf("expected count 1, got %v", e)
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	a := newWithClock(time.Minute, fixedClock(epoch))
	a.Record("tcp:8080")
	a.Record("tcp:8080")
	a.mu.Lock()
	c := a.bucket["tcp:8080"].Count
	a.mu.Unlock()
	if c != 2 {
		t.Fatalf("expected 2, got %d", c)
	}
}

func TestFlush_BeforeWindowReturnsNothing(t *testing.T) {
	a := newWithClock(time.Minute, fixedClock(epoch))
	a.Record("tcp:9090")
	out := a.Flush()
	if len(out) != 0 {
		t.Fatalf("expected empty flush, got %d events", len(out))
	}
}

func TestFlush_AfterWindowReturnsEvent(t *testing.T) {
	now := epoch
	a := newWithClock(time.Minute, func() time.Time { return now })
	a.Record("tcp:9090")
	now = epoch.Add(2 * time.Minute)
	out := a.Flush()
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
	if out[0].Key != "tcp:9090" {
		t.Errorf("unexpected key: %s", out[0].Key)
	}
}

func TestFlush_RemovesFromBucket(t *testing.T) {
	now := epoch
	a := newWithClock(time.Minute, func() time.Time { return now })
	a.Record("udp:53")
	now = epoch.Add(2 * time.Minute)
	a.Flush()
	a.mu.Lock()
	_, exists := a.bucket["udp:53"]
	a.mu.Unlock()
	if exists {
		t.Error("expected bucket entry to be removed after flush")
	}
}

func TestSummary_Format(t *testing.T) {
	e := Event{Key: "tcp:443", Count: 5, FirstSeen: epoch, LastSeen: epoch.Add(time.Minute)}
	s := Summary(e)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
