package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleEvent() Event {
	return Event{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Kind:      "new_binding",
		Port:      8080,
		Protocol:  "tcp",
		PID:       1234,
		Message:   "unexpected binding detected",
	}
}

func TestNewLogger_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "audit.log")
	l, err := NewLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("expected dir to exist: %v", err)
	}
}

func TestLog_And_ReadAll_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	l, _ := NewLogger(path)

	e := sampleEvent()
	if err := l.Log(e); err != nil {
		t.Fatalf("log error: %v", err)
	}

	events, err := ReadAll(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Port != e.Port {
		t.Errorf("port mismatch: got %d", events[0].Port)
	}
	if events[0].Kind != e.Kind {
		t.Errorf("kind mismatch: got %s", events[0].Kind)
	}
}

func TestReadAll_NonExistentFile(t *testing.T) {
	events, err := ReadAll("/nonexistent/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != nil {
		t.Errorf("expected nil events")
	}
}

func TestLog_AppendsMultipleEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	l, _ := NewLogger(path)

	for i := 0; i < 3; i++ {
		e := sampleEvent()
		e.Port = 8080 + i
		_ = l.Log(e)
	}

	events, err := ReadAll(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}

func TestLog_SetsTimestampIfZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	l, _ := NewLogger(path)

	e := Event{Kind: "test", Port: 9090, Protocol: "udp"}
	_ = l.Log(e)

	events, _ := ReadAll(path)
	if events[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}
