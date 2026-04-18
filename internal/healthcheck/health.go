package healthcheck

import (
	"sync"
	"time"
)

// Status represents the health state of a component.
type Status string

const (
	StatusOK      Status = "ok"
	StatusDegraded Status = "degraded"
	StatusUnknown  Status = "unknown"
)

// ComponentHealth holds the health info for a single component.
type ComponentHealth struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Message   string    `json:"message,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

// Report is the aggregated health of all components.
type Report struct {
	Overall    Status            `json:"overall"`
	Components []ComponentHealth `json:"components"`
	GeneratedAt time.Time        `json:"generated_at"`
}

// Checker is a function that returns a ComponentHealth.
type Checker func() ComponentHealth

// Monitor aggregates multiple health checkers.
type Monitor struct {
	mu       sync.RWMutex
	checkers map[string]Checker
	now      func() time.Time
}

// New returns a new Monitor.
func New() *Monitor {
	return &Monitor{
		checkers: make(map[string]Checker),
		now:      time.Now,
	}
}

// Register adds a named checker to the monitor.
func (m *Monitor) Register(name string, c Checker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkers[name] = c
}

// Check runs all checkers and returns an aggregated Report.
func (m *Monitor) Check() Report {
	m.mu.RLock()
	defer m.mu.RUnlock()

	report := Report{GeneratedAt: m.now(), Overall: StatusOK}
	for _, c := range m.checkers {
		ch := c()
		report.Components = append(report.Components, ch)
		if ch.Status == StatusDegraded {
			report.Overall = StatusDegraded
		} else if ch.Status == StatusUnknown && report.Overall != StatusDegraded {
			report.Overall = StatusUnknown
		}
	}
	return report
}
