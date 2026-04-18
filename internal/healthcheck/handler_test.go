package healthcheck

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestReporter_Write_ValidJSON(t *testing.T) {
	m := newTestMonitor()
	m.Register("scanner", StaticChecker("scanner", StatusOK, ""))

	var buf bytes.Buffer
	r := NewReporter(m, &buf)
	r.now = func() interface{} { return fixedNow }.(func() interface{})

	if err := r.Write(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report Report
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if report.Overall != StatusOK {
		t.Errorf("expected OK, got %s", report.Overall)
	}
	if len(report.Components) != 1 {
		t.Errorf("expected 1 component, got %d", len(report.Components))
	}
}

func TestReporter_Write_DegradedReflectedInJSON(t *testing.T) {
	m := newTestMonitor()
	m.Register("alerter", StaticChecker("alerter", StatusDegraded, "out of memory"))

	var buf bytes.Buffer
	r := NewReporter(m, &buf)

	if err := r.Write(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report Report
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if report.Overall != StatusDegraded {
		t.Errorf("expected Degraded, got %s", report.Overall)
	}
}
