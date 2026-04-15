package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Snapshot holds a point-in-time capture of observed port bindings.
type Snapshot struct {
	CapturedAt time.Time              `json:"captured_at"`
	Entries    []portscanner.PortEntry `json:"entries"`
}

// Store persists snapshots to disk and retrieves the last saved one.
type Store struct {
	path string
}

// NewStore creates a Store that reads/writes to the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the current snapshot to disk, creating directories as needed.
func (s *Store) Save(snap Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(snap)
}

// Load reads the last persisted snapshot from disk.
// Returns an empty Snapshot and no error when the file does not exist yet.
func (s *Store) Load() (Snapshot, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
