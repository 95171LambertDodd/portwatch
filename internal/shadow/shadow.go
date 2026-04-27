// Package shadow provides a shadow-mode evaluator that runs detection logic
// without emitting real alerts, logging what would have fired instead.
package shadow

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Event represents a detection that would have fired in live mode.
type Event struct {
	Timestamp time.Time
	Port      uint16
	Protocol  string
	Process   string
	Reason    string
}

// Evaluator runs checks in shadow mode, collecting would-be alerts.
type Evaluator struct {
	mu     sync.Mutex
	events []Event
	out    io.Writer
	clock  func() time.Time
}

// New creates a shadow Evaluator that writes shadow-mode notices to out.
func New(out io.Writer) (*Evaluator, error) {
	if out == nil {
		return nil, fmt.Errorf("shadow: output writer must not be nil")
	}
	return &Evaluator{
		out:   out,
		clock: time.Now,
	}, nil
}

func newWithClock(out io.Writer, clock func() time.Time) *Evaluator {
	return &Evaluator{out: out, clock: clock}
}

// Evaluate checks whether the entry would trigger an alert given the known
// set of baseline ports. It records and logs the event without side-effects.
func (e *Evaluator) Evaluate(entry portscanner.PortEntry, knownPorts map[uint16]bool) {
	if knownPorts[entry.Port] {
		return
	}
	ev := Event{
		Timestamp: e.clock(),
		Port:      entry.Port,
		Protocol:  entry.Protocol,
		Process:   entry.Process,
		Reason:    "unrecognised port binding",
	}
	e.mu.Lock()
	e.events = append(e.events, ev)
	e.mu.Unlock()
	fmt.Fprintf(e.out, "[shadow] %s port=%d proto=%s process=%q reason=%s\n",
		ev.Timestamp.Format(time.RFC3339), ev.Port, ev.Protocol, ev.Process, ev.Reason)
}

// Events returns a snapshot of all collected shadow events.
func (e *Evaluator) Events() []Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]Event, len(e.events))
	copy(out, e.events)
	return out
}

// Reset clears all collected shadow events.
func (e *Evaluator) Reset() {
	e.mu.Lock()
	e.events = e.events[:0]
	e.mu.Unlock()
}
