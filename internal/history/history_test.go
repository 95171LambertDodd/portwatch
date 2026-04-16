package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleEvent(kind string) Event {
	return Event{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Proto:     "tcp",
		Addr:      "0.0.0.0",
		Port:      8080,
		PID:       1234,
		Process:   "nginx",
		Kind:      kind,
	}
}

func TestNewRecorder_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "history.jsonl")
	r, err := NewRecorder(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil recorder")
	}
}

func TestRecord_And_ReadAll_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.jsonl")

	r, err := NewRecorder(path)
	if err != nil {
		t.Fatalf("NewRecorder: %v", err)
	}

	events := []Event{sampleEvent("new"), sampleEvent("gone")}
	for _, e := range events {
		if err := r.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	loaded, err := ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 events, got %d", len(loaded))
	}
	if loaded[0].Kind != "new" || loaded[1].Kind != "gone" {
		t.Errorf("unexpected kinds: %v %v", loaded[0].Kind, loaded[1].Kind)
	}
}

func TestReadAll_NonExistentFile(t *testing.T) {
	events, err := ReadAll("/tmp/portwatch_no_such_file_xyz.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if events != nil {
		t.Errorf("expected nil events, got %v", events)
	}
}

func TestRecord_AppendsMultipleTimes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.jsonl")

	r, _ := NewRecorder(path)
	for i := 0; i < 5; i++ {
		_ = r.Record(sampleEvent("new"))
	}

	data, _ := os.ReadFile(path)
	lines := splitLines(data)
	count := 0
	for _, l := range lines {
		if len(l) > 0 {
			count++
		}
	}
	if count != 5 {
		t.Errorf("expected 5 lines, got %d", count)
	}
}
