package retention

import (
	"testing"
)

func TestNewSink_NilManager_ReturnsError(t *testing.T) {
	_, err := NewSink(nil)
	if err == nil {
		t.Fatal("expected error for nil manager")
	}
}

func TestNewSink_ValidManager_ReturnsNonNil(t *testing.T) {
	m, _ := newWithClock(Policy{MaxCount: 5}, fixedClock(epoch))
	s, err := NewSink(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sink")
	}
}

func TestSink_Write_RecordsEntry(t *testing.T) {
	m, _ := newWithClock(Policy{MaxCount: 10}, fixedClock(epoch))
	s, _ := NewSink(m)

	if err := s.Write(makeEntry(8080)); err != nil {
		t.Fatalf("Write returned unexpected error: %v", err)
	}

	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry in snapshot, got %d", len(snap))
	}
	if snap[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", snap[0].Port)
	}
}

func TestSink_Snapshot_ReflectsPrune(t *testing.T) {
	m, _ := newWithClock(Policy{MaxCount: 2}, fixedClock(epoch))
	s, _ := NewSink(m)

	_ = s.Write(makeEntry(1))
	_ = s.Write(makeEntry(2))
	_ = s.Write(makeEntry(3))

	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries after prune, got %d", len(snap))
	}
	if snap[0].Port != 2 {
		t.Errorf("expected first retained port to be 2, got %d", snap[0].Port)
	}
}
