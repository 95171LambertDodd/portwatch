package shadow

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

func fixedClock() time.Time {
	return time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
}

func makeEntry(port uint16, proto, process string) portscanner.PortEntry {
	return portscanner.PortEntry{Port: port, Protocol: proto, Process: process}
}

func TestNew_NilWriter_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestNew_ValidWriter_ReturnsNonNil(t *testing.T) {
	ev, err := New(&bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev == nil {
		t.Fatal("expected non-nil evaluator")
	}
}

func TestEvaluate_KnownPort_NoEvent(t *testing.T) {
	var buf bytes.Buffer
	ev := newWithClock(&buf, fixedClock)
	known := map[uint16]bool{8080: true}
	ev.Evaluate(makeEntry(8080, "tcp", "nginx"), known)
	if len(ev.Events()) != 0 {
		t.Fatalf("expected 0 events, got %d", len(ev.Events()))
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %s", buf.String())
	}
}

func TestEvaluate_UnknownPort_RecordsEvent(t *testing.T) {
	var buf bytes.Buffer
	ev := newWithClock(&buf, fixedClock)
	known := map[uint16]bool{}
	ev.Evaluate(makeEntry(9999, "tcp", "mystery"), known)
	events := ev.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", events[0].Port)
	}
	if events[0].Process != "mystery" {
		t.Errorf("expected process 'mystery', got %q", events[0].Process)
	}
}

func TestEvaluate_WritesToOutput(t *testing.T) {
	var buf bytes.Buffer
	ev := newWithClock(&buf, fixedClock)
	ev.Evaluate(makeEntry(1234, "udp", "rogue"), map[uint16]bool{})
	out := buf.String()
	if !strings.Contains(out, "[shadow]") {
		t.Errorf("expected [shadow] prefix in output, got: %s", out)
	}
	if !strings.Contains(out, "port=1234") {
		t.Errorf("expected port=1234 in output, got: %s", out)
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	var buf bytes.Buffer
	ev := newWithClock(&buf, fixedClock)
	ev.Evaluate(makeEntry(5000, "tcp", "svc"), map[uint16]bool{})
	ev.Reset()
	if len(ev.Events()) != 0 {
		t.Fatalf("expected 0 events after reset, got %d", len(ev.Events()))
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	var buf bytes.Buffer
	ev := newWithClock(&buf, fixedClock)
	ev.Evaluate(makeEntry(7000, "tcp", "svc"), map[uint16]bool{})
	snap := ev.Events()
	snap[0].Port = 0
	original := ev.Events()
	if original[0].Port != 7000 {
		t.Errorf("Events() should return a copy; original was mutated")
	}
}
