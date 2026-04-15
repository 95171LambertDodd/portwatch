package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/snapshot"
)

func TestStore_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	store := snapshot.NewStore(path)

	original := snapshot.Snapshot{
		CapturedAt: time.Now().UTC().Truncate(time.Second),
		Entries: []portscanner.PortEntry{
			{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 8080, PID: 42, State: "LISTEN"},
		},
	}

	if err := store.Save(original); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(loaded.Entries) != len(original.Entries) {
		t.Errorf("expected %d entries, got %d", len(original.Entries), len(loaded.Entries))
	}
	if loaded.Entries[0].LocalPort != 8080 {
		t.Errorf("expected port 8080, got %d", loaded.Entries[0].LocalPort)
	}
}

func TestStore_Load_NonExistentFile(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(filepath.Join(dir, "missing.json"))

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(snap.Entries) != 0 {
		t.Errorf("expected empty entries for len(snap.Entries))
	}
}

func TestStore_Save_CreatesIntermediateDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "snap.json")
	storen
	if err := store.Save(snapshot.Snapshot{CapturedAt: time.Now()}); err != nil {
		t.Fatalf("Save() should create dirs, got error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist at %s: %v", path, err)
	}
}

func TestStore_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	store := snapshot.NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
