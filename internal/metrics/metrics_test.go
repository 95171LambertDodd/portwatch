package metrics

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_ReturnsNonNil(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil Collector")
	}
}

func TestRecordScan_IncrementsCount(t *testing.T) {
	c := New()
	c.RecordScan()
	c.RecordScan()
	if got := c.Snapshot().ScanCount; got != 2 {
		t.Fatalf("expected ScanCount=2, got %d", got)
	}
}

func TestRecordAlert_IncrementsCount(t *testing.T) {
	c := New()
	c.RecordAlert()
	if got := c.Snapshot().AlertCount; got != 1 {
		t.Fatalf("expected AlertCount=1, got %d", got)
	}
}

func TestRecordViolation_IncrementsCount(t *testing.T) {
	c := New()
	c.RecordViolation()
	c.RecordViolation()
	c.RecordViolation()
	if got := c.Snapshot().ViolationCount; got != 3 {
		t.Fatalf("expected ViolationCount=3, got %d", got)
	}
}

func TestRecordScan_SetsLastScanAt(t *testing.T) {
	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	c := New()
	c.now = fixedNow(ts)
	c.RecordScan()
	if got := c.Snapshot().LastScanAt; !got.Equal(ts) {
		t.Fatalf("expected LastScanAt=%v, got %v", ts, got)
	}
}

func TestRecordAlert_SetsLastAlertAt(t *testing.T) {
	ts := time.Date(2024, 6, 1, 8, 30, 0, 0, time.UTC)
	c := New()
	c.now = fixedNow(ts)
	c.RecordAlert()
	if got := c.Snapshot().LastAlertAt; !got.Equal(ts) {
		t.Fatalf("expected LastAlertAt=%v, got %v", ts, got)
	}
}

func TestReset_ZeroesAllFields(t *testing.T) {
	c := New()
	c.RecordScan()
	c.RecordAlert()
	c.RecordViolation()
	c.Reset()
	s := c.Snapshot()
	if s.ScanCount != 0 || s.AlertCount != 0 || s.ViolationCount != 0 {
		t.Fatalf("expected all zeros after Reset, got %+v", s)
	}
	if !s.LastScanAt.IsZero() || !s.LastAlertAt.IsZero() {
		t.Fatal("expected zero timestamps after Reset")
	}
}
