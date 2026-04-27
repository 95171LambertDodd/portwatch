package shadow

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Sink is a notify-compatible sink that routes entries through the shadow
// Evaluator instead of emitting real alerts.
type Sink struct {
	eval      *Evaluator
	known     map[uint16]bool
	jsonOut   io.Writer
}

// NewSink creates a Sink backed by the given Evaluator. knownPorts is the
// set of baseline ports that should not trigger shadow events. jsonOut
// receives a JSON summary when Flush is called; it may be nil.
func NewSink(eval *Evaluator, knownPorts map[uint16]bool, jsonOut io.Writer) (*Sink, error) {
	if eval == nil {
		return nil, fmt.Errorf("shadow: evaluator must not be nil")
	}
	if knownPorts == nil {
		knownPorts = map[uint16]bool{}
	}
	return &Sink{eval: eval, known: knownPorts, jsonOut: jsonOut}, nil
}

// Write evaluates a single PortEntry in shadow mode.
func (s *Sink) Write(entry portscanner.PortEntry) error {
	s.eval.Evaluate(entry, s.known)
	return nil
}

// Flush writes a JSON summary of collected shadow events to jsonOut.
// It is a no-op when jsonOut is nil.
func (s *Sink) Flush() error {
	if s.jsonOut == nil {
		return nil
	}
	type summary struct {
		FlushedAt   time.Time `json:"flushed_at"`
		EventCount  int       `json:"event_count"`
		Events      []Event   `json:"events"`
	}
	events := s.eval.Events()
	enc := json.NewEncoder(s.jsonOut)
	return enc.Encode(summary{
		FlushedAt:  time.Now(),
		EventCount: len(events),
		Events:     events,
	})
}
