package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents a single port binding event recorded in history.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Proto     string    `json:"proto"`
	Addr      string    `json:"addr"`
	Port      uint16    `json:"port"`
	PID       int       `json:"pid"`
	Process   string    `json:"process"`
	Kind      string    `json:"kind"` // "new" | "gone"
}

// Recorder appends port events to a JSON-lines file.
type Recorder struct {
	mu   sync.Mutex
	path string
}

// NewRecorder creates a Recorder that writes to the given file path.
func NewRecorder(path string) (*Recorder, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return &Recorder{path: path}, nil
}

// Record appends an event to the history file.
func (r *Recorder) Record(e Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	line, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = f.Write(append(line, '\n'))
	return err
}

// ReadAll loads all events from the history file.
func ReadAll(path string) ([]Event, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var events []Event
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Event
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
