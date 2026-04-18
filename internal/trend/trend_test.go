package trend

import (
	"testing"
	"time"
)

var baseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_ReturnsNonNil(t *testing.T) {
	tr := New(time.Minute)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestRate_NoData_ReturnsZero(t *testing.T) {
	tr := newWithClock(time.Minute, fixedClock(baseTime))
	if r := tr.Rate("tcp:8080"); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestRate_SinglePoint_ReturnsZero(t *testing.T) {
	tr := newWithClock(time.Minute, fixedClock(baseTime))
	tr.Record("tcp:8080", 5)
	if r := tr.Rate("tcp:8080"); r != 0 {
		t.Fatalf("expected 0 with single point, got %f", r)
	}
}

func TestRate_MultiplePoints_ReturnsAverage(t *testing.T) {
	now := baseTime
	calls := []time.Time{now, now.Add(10 * time.Second), now.Add(20 * time.Second)}
	idx := 0
	clock := func() time.Time {
		t := calls[idx]
		if idx < len(calls)-1 {
			idx++
		}
		return t
	}
	tr := newWithClock(time.Minute, clock)
	tr.Record("tcp:80", 2)
	tr.Record("tcp:80", 4)
	tr.Record("tcp:80", 6)
	r := tr.Rate("tcp:80")
	if r != 4.0 {
		t.Fatalf("expected 4.0, got %f", r)
	}
}

func TestRate_PrunesOldPoints(t *testing.T) {
	now := baseTime
	tr := newWithClock(time.Minute, fixedClock(now))
	// Record an old point manually by manipulating time.
	tr.points["tcp:443"] = []Point{
		{At: now.Add(-2 * time.Minute), Count: 100},
		{At: now.Add(-90 * time.Second), Count: 50},
	}
	// Both points are outside the 1-minute window.
	if r := tr.Rate("tcp:443"); r != 0 {
		t.Fatalf("expected 0 after pruning, got %f", r)
	}
}

func TestRecord_DifferentKeys_Independent(t *testing.T) {
	tr := newWithClock(time.Minute, fixedClock(baseTime))
	tr.Record("tcp:80", 10)
	tr.Record("udp:53", 20)
	if r := tr.Rate("tcp:80"); r != 0 {
		t.Fatal("single point should return 0")
	}
	if _, ok := tr.points["udp:53"]; !ok {
		t.Fatal("expected udp:53 key to exist")
	}
}
