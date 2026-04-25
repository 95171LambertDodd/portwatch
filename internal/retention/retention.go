// Package retention provides a policy-based retention manager that prunes
// old port-event records beyond a configurable age or count limit.
package retention

import (
	"errors"
	"sync"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Policy describes how long or how many events to retain.
type Policy struct {
	// MaxAge is the maximum age of an event before it is pruned.
	// A zero value disables age-based pruning.
	MaxAge time.Duration

	// MaxCount is the maximum number of events to retain (oldest pruned first).
	// A zero value disables count-based pruning.
	MaxCount int
}

// Manager applies a retention Policy to a slice of port entries.
type Manager struct {
	policy Policy
	clock  func() time.Time
	mu     sync.Mutex
	events []timedEntry
}

type timedEntry struct {
	entry     portscanner.PortEntry
	recorded  time.Time
}

// New creates a Manager with the given Policy.
// Returns an error if the policy is entirely unconstrained.
func New(p Policy) (*Manager, error) {
	return newWithClock(p, time.Now)
}

func newWithClock(p Policy, clock func() time.Time) (*Manager, error) {
	if p.MaxAge == 0 && p.MaxCount == 0 {
		return nil, errors.New("retention: policy must set MaxAge, MaxCount, or both")
	}
	return &Manager{policy: p, clock: clock}, nil
}

// Record adds a PortEntry to the managed set, then applies the retention policy.
func (m *Manager) Record(e portscanner.PortEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, timedEntry{entry: e, recorded: m.clock()})
	m.prune()
}

// Entries returns a copy of all currently retained entries.
func (m *Manager) Entries() []portscanner.PortEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]portscanner.PortEntry, len(m.events))
	for i, te := range m.events {
		out[i] = te.entry
	}
	return out
}

// prune must be called with m.mu held.
func (m *Manager) prune() {
	if m.policy.MaxAge > 0 {
		cutoff := m.clock().Add(-m.policy.MaxAge)
		start := 0
		for start < len(m.events) && m.events[start].recorded.Before(cutoff) {
			start++
		}
		m.events = m.events[start:]
	}
	if m.policy.MaxCount > 0 && len(m.events) > m.policy.MaxCount {
		m.events = m.events[len(m.events)-m.policy.MaxCount:]
	}
}
