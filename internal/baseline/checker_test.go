package baseline

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portscanner"
)

func makePortEntry(proto, addr string, port uint16, pid int) portscanner.PortEntry {
	return portscanner.PortEntry{
		Proto:        proto,
		LocalAddress: addr,
		LocalPort:    port,
		PID:          pid,
	}
}

func baselineWithEntries(entries ...Entry) *Baseline {
	return &Baseline{Entries: entries}
}

func TestChecker_NoViolationsWhenAllKnown(t *testing.T) {
	b := baselineWithEntries(
		Entry{Proto: "tcp", Address: "0.0.0.0", Port: 80},
		Entry{Proto: "tcp", Address: "0.0.0.0", Port: 443},
	)
	c := NewChecker(b, nil)
	entries := []portscanner.PortEntry{
		makePortEntry("tcp", "0.0.0.0", 80, 100),
		makePortEntry("tcp", "0.0.0.0", 443, 101),
	}
	violations := c.Check(entries)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestChecker_ViolationForUnknownPort(t *testing.T) {
	b := baselineWithEntries(Entry{Proto: "tcp", Address: "0.0.0.0", Port: 80})
	var buf bytes.Buffer
	c := NewChecker(b, &buf)
	entries := []portscanner.PortEntry{
		makePortEntry("tcp", "0.0.0.0", 80, 100),
		makePortEntry("tcp", "0.0.0.0", 9999, 200),
	}
	violations := c.Check(entries)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Entry.LocalPort != 9999 {
		t.Errorf("expected port 9999, got %d", violations[0].Entry.LocalPort)
	}
	if !strings.Contains(buf.String(), "VIOLATION") {
		t.Error("expected VIOLATION in output")
	}
}

func TestChecker_EmptyBaselineAllViolations(t *testing.T) {
	b := &Baseline{}
	c := NewChecker(b, nil)
	entries := []portscanner.PortEntry{
		makePortEntry("tcp", "127.0.0.1", 22, 1),
		makePortEntry("udp", "0.0.0.0", 53, 2),
	}
	violations := c.Check(entries)
	if len(violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(violations))
	}
}

func TestChecker_EmptyEntries(t *testing.T) {
	b := baselineWithEntries(Entry{Proto: "tcp", Address: "0.0.0.0", Port: 80})
	c := NewChecker(b, nil)
	violations := c.Check(nil)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{
		Entry:  makePortEntry("tcp", "0.0.0.0", 8080, 42),
		Reason: "not in baseline",
	}
	s := v.String()
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port in string, got: %s", s)
	}
	if !strings.Contains(s, "not in baseline") {
		t.Errorf("expected reason in string, got: %s", s)
	}
}
