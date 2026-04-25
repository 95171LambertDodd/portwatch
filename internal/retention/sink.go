package retention

import (
	"errors"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Sink is a notify-compatible sink that feeds entries into a retention Manager.
// It satisfies the same Write(portscanner.PortEntry) error interface used
// elsewhere in the pipeline.
type Sink struct {
	manager  *Manager
	maxAge   time.Duration
	maxCount int
}

// NewSink creates a Sink backed by the provided Manager.
func NewSink(m *Manager) (*Sink, error) {
	if m == nil {
		return nil, errors.New("retention: manager must not be nil")
	}
	return &Sink{manager: m}, nil
}

// Write records the entry into the retention Manager.
func (s *Sink) Write(e portscanner.PortEntry) error {
	s.manager.Record(e)
	return nil
}

// Snapshot returns all currently retained entries.
func (s *Sink) Snapshot() []portscanner.PortEntry {
	return s.manager.Entries()
}
