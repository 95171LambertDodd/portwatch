package replay_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/replay"
	"github.com/user/portwatch/internal/portscanner"
)

// sampleEntry returns a minimal PortEntry for use in tests.
func sampleEntry(port uint16, proto string) portscanner.PortEntry {
	return portscanner.PortEntry{
		Port:     port,
		Protocol: proto,
		PID:      1234,
	}
}

func TestNew_ReturnsNonNil(t *testing.T) {
	r := replay.New(10)
	if r == nil {
		t.Fatal("expected non-nil Replayer")
	}
}

func TestRecord_StoresEvent(t *testing.T) {
	r := replay.New(10)
	entry := sampleEntry(8080, "tcp")

	r.Record(entry, "new_binding", time.Now())

	events := r.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Entry.Port != 8080 {
		t.Errorf("expected port 8080, got %d", events[0].Entry.Port)
	}
	if events[0].Kind != "new_binding" {
		t.Errorf("expected kind 'new_binding', got %q", events[0].Kind)
	}
}

func TestRecord_RespectsCapacity(t *testing.T) {
	cap := 3
	r := replay.New(cap)

	for i := 0; i < 5; i++ {
		r.Record(sampleEntry(uint16(8000+i), "tcp"), "new_binding", time.Now())
	}

	events := r.Events()
	if len(events) != cap {
		t.Fatalf("expected %d events (capacity), got %d", cap, len(events))
	}
	// The oldest entries should have been evicted; last 3 ports should be 8002, 8003, 8004.
	if events[0].Entry.Port != 8002 {
		t.Errorf("expected oldest retained port 8002, got %d", events[0].Entry.Port)
	}
	if events[2].Entry.Port != 8004 {
		t.Errorf("expected newest port 8004, got %d", events[2].Entry.Port)
	}
}

func TestEvents_ReturnsCopyInOrder(t *testing.T) {
	r := replay.New(10)
	now := time.Now()

	r.Record(sampleEntry(9001, "tcp"), "new_binding", now)
	r.Record(sampleEntry(9002, "udp"), "resolved", now.Add(time.Second))

	events := r.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Entry.Port != 9001 {
		t.Errorf("expected first event port 9001, got %d", events[0].Entry.Port)
	}
	if events[1].Entry.Port != 9002 {
		t.Errorf("expected second event port 9002, got %d", events[1].Entry.Port)
	}
}

func TestEvents_MutatingReturnDoesNotAffectInternal(t *testing.T) {
	r := replay.New(10)
	r.Record(sampleEntry(443, "tcp"), "new_binding", time.Now())

	events := r.Events()
	events[0].Kind = "tampered"

	// Re-fetch and confirm original is unchanged.
	fresh := r.Events()
	if fresh[0].Kind == "tampered" {
		t.Error("mutating returned slice affected internal state")
	}
}

func TestClear_RemovesAllEvents(t *testing.T) {
	r := replay.New(10)
	r.Record(sampleEntry(80, "tcp"), "new_binding", time.Now())
	r.Record(sampleEntry(443, "tcp"), "new_binding", time.Now())

	r.Clear()

	events := r.Events()
	if len(events) != 0 {
		t.Fatalf("expected 0 events after Clear, got %d", len(events))
	}
}

func TestRecord_TimestampPreserved(t *testing.T) {
	r := replay.New(10)
	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	r.Record(sampleEntry(8080, "tcp"), "new_binding", ts)

	events := r.Events()
	if !events[0].Timestamp.Equal(ts) {
		t.Errorf("expected timestamp %v, got %v", ts, events[0].Timestamp)
	}
}
