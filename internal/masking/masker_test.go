package masking_test

import (
	"testing"

	"github.com/user/portwatch/internal/masking"
	"github.com/user/portwatch/internal/portscanner"
)

func makeEntry() portscanner.PortEntry {
	return portscanner.PortEntry{
		LocalAddr: "0.0.0.0:8080",
		Protocol:  "tcp",
		PID:       1234,
		Process:   "nginx",
	}
}

func TestNew_ValidRules_ReturnsNonNil(t *testing.T) {
	m, err := masking.New([]masking.Rule{{Field: "pid", Replacement: "[redacted]"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil masker")
	}
}

func TestNew_InvalidField_ReturnsError(t *testing.T) {
	_, err := masking.New([]masking.Rule{{Field: "username", Replacement: "x"}})
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestNew_EmptyReplacement_ReturnsError(t *testing.T) {
	_, err := masking.New([]masking.Rule{{Field: "process", Replacement: ""}})
	if err == nil {
		t.Fatal("expected error for empty replacement")
	}
}

func TestApply_MasksProcess(t *testing.T) {
	m, _ := masking.New([]masking.Rule{{Field: "process", Replacement: "[hidden]"}})
	out := m.Apply(makeEntry())
	if out.Process != "[hidden]" {
		t.Errorf("expected [hidden], got %q", out.Process)
	}
	if out.PID != 1234 {
		t.Errorf("PID should be unchanged, got %d", out.PID)
	}
}

func TestApply_MasksPID(t *testing.T) {
	m, _ := masking.New([]masking.Rule{{Field: "pid", Replacement: "0"}})
	out := m.Apply(makeEntry())
	if out.PID != 0 {
		t.Errorf("expected PID 0, got %d", out.PID)
	}
}

func TestApply_MasksAddr(t *testing.T) {
	m, _ := masking.New([]masking.Rule{{Field: "addr", Replacement: "0.0.0.0:0"}})
	out := m.Apply(makeEntry())
	if out.LocalAddr != "0.0.0.0:0" {
		t.Errorf("expected masked addr, got %q", out.LocalAddr)
	}
}

func TestApplyAll_ReturnsMaskedSlice(t *testing.T) {
	m, _ := masking.New([]masking.Rule{{Field: "process", Replacement: "[redacted]"}})
	entries := []portscanner.PortEntry{makeEntry(), makeEntry()}
	out := m.ApplyAll(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Process != "[redacted]" {
			t.Errorf("expected process redacted, got %q", e.Process)
		}
	}
}

func TestApplyAll_OriginalUnchanged(t *testing.T) {
	m, _ := masking.New([]masking.Rule{{Field: "process", Replacement: "[redacted]"}})
	original := []portscanner.PortEntry{makeEntry()}
	m.ApplyAll(original)
	if original[0].Process != "nginx" {
		t.Errorf("original slice must not be mutated")
	}
}
