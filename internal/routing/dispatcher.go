package routing

import (
	"errors"
	"fmt"
	"io"

	"github.com/user/portwatch/internal/portscanner"
)

// Writer is the interface for a named output sink.
type Writer interface {
	Write(entry portscanner.PortEntry) error
}

// Dispatcher combines a Router with a map of named Writers and dispatches
// each entry to the appropriate sink.
type Dispatcher struct {
	router *Router
	sinks  map[string]Writer
	dead   io.Writer // receives unrouted entries as JSON-ish log lines; may be nil
}

// NewDispatcher creates a Dispatcher. sinks must contain every destination
// referenced by the router's rules.
func NewDispatcher(router *Router, sinks map[string]Writer, dead io.Writer) (*Dispatcher, error) {
	if router == nil {
		return nil, errors.New("router must not be nil")
	}
	if len(sinks) == 0 {
		return nil, errors.New("at least one sink is required")
	}
	for _, rule := range router.rules {
		if _, ok := sinks[rule.Destination]; !ok {
			return nil, fmt.Errorf("no sink registered for destination %q", rule.Destination)
		}
	}
	return &Dispatcher{router: router, sinks: sinks, dead: dead}, nil
}

// Dispatch routes entry to the matching sink and calls Write.
// If no rule matches and a dead-letter writer is configured the entry is
// logged there; otherwise the entry is silently dropped.
func (d *Dispatcher) Dispatch(e portscanner.PortEntry) error {
	dest := d.router.Route(e)
	if dest == "" {
		if d.dead != nil {
			fmt.Fprintf(d.dead, "[routing] unrouted entry: proto=%s port=%d pid=%d\n",
				e.Protocol, e.Port, e.PID)
		}
		return nil
	}
	sink, ok := d.sinks[dest]
	if !ok {
		return fmt.Errorf("sink %q not found", dest)
	}
	return sink.Write(e)
}
