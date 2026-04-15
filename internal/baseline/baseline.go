package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a known/expected port binding in the baseline.
type Entry struct {
	Proto   string `json:"proto"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
	PID     int    `json:"pid"`
	Comment string `json:"comment,omitempty"`
}

// Baseline holds the set of expected port bindings.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

// Manager handles loading, saving, and querying the baseline file.
type Manager struct {
	path string
}

// NewManager creates a new Manager for the given file path.
func NewManager(path string) *Manager {
	return &Manager{path: path}
}

// Load reads the baseline from disk. Returns an empty Baseline if the file does not exist.
func (m *Manager) Load() (*Baseline, error) {
	data, err := os.ReadFile(m.path)
	if os.IsNotExist(err) {
		return &Baseline{CreatedAt: time.Now()}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", m.path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: parse %s: %w", m.path, err)
	}
	return &b, nil
}

// Save writes the baseline to disk, creating intermediate directories as needed.
func (m *Manager) Save(b *Baseline) error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o755); err != nil {
		return fmt.Errorf("baseline: mkdir %s: %w", filepath.Dir(m.path), err)
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(m.path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", m.path, err)
	}
	return nil
}

// Contains reports whether the given proto/address/port combination is in the baseline.
func (b *Baseline) Contains(proto, address string, port uint16) bool {
	for _, e := range b.Entries {
		if e.Proto == proto && e.Address == address && e.Port == port {
			return true
		}
	}
	return false
}
