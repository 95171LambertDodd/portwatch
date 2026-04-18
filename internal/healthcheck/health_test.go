package healthcheck

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func newTestMonitor() *Monitor {
	m := New()
	m.now = func() time.Time { return fixedNow }
	return m
}

func TestNew_ReturnsNonNil(t *testing.T) {
	if New() == nil {
		t.Fatal("expected non-nil monitor")
	}
}

func TestCheck_NoCheckers_ReturnsOK(t *testing.T) {
	m := newTestMonitor()
	r := m.Check()
	if r.Overall != StatusOK {
		t.Errorf("expected OK, got %s", r.Overall)
	}
	if len(r.Components) != 0 {
		t.Errorf("expected no components")
	}
}

func TestCheck_AllOK(t *testing.T) {
	m := newTestMonitor()
	m.Register("scanner", func() ComponentHealth {
		return ComponentHealth{Name: "scanner", Status: StatusOK}
	})
	r := m.Check()
	if r.Overall != StatusOK {
		t.Errorf("expected OK, got %s", r.Overall)
	}
}

func TestCheck_OneDegraded_OverallDegraded(t *testing.T) {
	m := newTestMonitor()
	m.Register("scanner", func() ComponentHealth {
		return ComponentHealth{Name: "scanner", Status: StatusOK}
	})
	m.Register("alerter", func() ComponentHealth {
		return ComponentHealth{Name: "alerter", Status: StatusDegraded, Message: "disk full"}
	})
	r := m.Check()
	if r.Overall != StatusDegraded {
		t.Errorf("expected Degraded, got %s", r.Overall)
	}
}

func TestCheck_UnknownWithNoOtherIssues(t *testing.T) {
	m := newTestMonitor()
	m.Register("notifier", func() ComponentHealth {
		return ComponentHealth{Name: "notifier", Status: StatusUnknown}
	})
	r := m.Check()
	if r.Overall != StatusUnknown {
		t.Errorf("expected Unknown, got %s", r.Overall)
	}
}

func TestCheck_DegradedTakesPriorityOverUnknown(t *testing.T) {
	m := newTestMonitor()
	m.Register("a", func() ComponentHealth {
		return ComponentHealth{Name: "a", Status: StatusUnknown}
	})
	m.Register("b", func() ComponentHealth {
		return ComponentHealth{Name: "b", Status: StatusDegraded}
	})
	r := m.Check()
	if r.Overall != StatusDegraded {
		t.Errorf("expected Degraded, got %s", r.Overall)
	}
}
