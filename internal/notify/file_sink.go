package notify

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileSink appends JSON-encoded notifications to a log file.
type FileSink struct {
	path string
}

// NewFileSink returns a FileSink that appends to the file at path.
// The file is created if it does not exist.
func NewFileSink(path string) *FileSink {
	return &FileSink{path: path}
}

// Send marshals msg as JSON and appends it to the log file.
func (s *FileSink) Send(msg Message) error {
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("file_sink: open %q: %w", s.path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(msg); err != nil {
		return fmt.Errorf("file_sink: encode: %w", err)
	}
	return nil
}
