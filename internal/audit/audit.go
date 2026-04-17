package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      string    `json:"kind"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	PID       int       `json:"pid,omitempty"`
	Message   string    `json:"message"`
}

// Logger writes audit events to a file in JSONL format.
type Logger struct {
	path string
}

// NewLogger creates a new Logger that writes to the given file path.
func NewLogger(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("audit: create dir: %w", err)
	}
	return &Logger{path: path}, nil
}

// Log appends an event to the audit log.
func (l *Logger) Log(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode event: %w", err)
	}
	return nil
}

// ReadAll reads all audit events from the log file.
func ReadAll(path string) ([]Event, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read file: %w", err)
	}
	var events []Event
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Event
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse line: %w", err)
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
