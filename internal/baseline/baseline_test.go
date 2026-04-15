package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager_ReturnsNonNil(t *testing.T) {
	m := NewManager("/tmp/baseline.json")
	if m == nil {
		t.Fatal("expected non-nil Manager")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	m := NewManager("/tmp/portwatch_nonexistent_baseline.json")
	b, err := m.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Baseline")
	}
	if len(b.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(b.Entries))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	m := NewManager(path)

	b := &Baseline{
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		Entries: []Entry{
			{Proto: "tcp", Address: "0.0.0.0", Port: 8080, PID: 1234, Comment: "web"},
			{Proto: "udp", Address: "127.0.0.1", Port: 53, PID: 5678},
		},
	}
	if err := m.Save(b); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := m.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", loaded.Entries[0].Port)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := NewManager(path)
	_, err := m.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSave_CreatesIntermediateDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "baseline.json")
	m := NewManager(path)
	b := &Baseline{CreatedAt: time.Now()}
	if err := m.Save(b); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestBaseline_Contains(t *testing.T) {
	b := &Baseline{
		Entries: []Entry{
			{Proto: "tcp", Address: "0.0.0.0", Port: 443},
		},
	}
	if !b.Contains("tcp", "0.0.0.0", 443) {
		t.Error("expected Contains to return true")
	}
	if b.Contains("tcp", "0.0.0.0", 80) {
		t.Error("expected Contains to return false for port 80")
	}
	if b.Contains("udp", "0.0.0.0", 443) {
		t.Error("expected Contains to return false for udp")
	}
}

func TestBaseline_ContainsEmpty(t *testing.T) {
	b := &Baseline{}
	if b.Contains("tcp", "127.0.0.1", 22) {
		t.Error("expected false on empty baseline")
	}
}

// Ensure Entry serialises cleanly (omitempty on Comment).
func TestEntry_JSONOmitEmpty(t *testing.T) {
	e := Entry{Proto: "tcp", Address: "0.0.0.0", Port: 80, PID: 1}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "" {
		t.Fatal("empty JSON")
	}
	// Comment field should be absent
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["comment"]; ok {
		t.Error("comment should be omitted when empty")
	}
}
