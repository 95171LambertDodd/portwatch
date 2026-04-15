package alerting

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portscanner"
)

func makeEntry(addr string, port uint16, pid int, proto string) portscanner.PortEntry {
	return portscanner.PortEntry{
		LocalAddr: addr,
		LocalPort: port,
		PID:       pid,
		Protocol:  proto,
	}
}

func TestDiff_NewBindingDetected(t *testing.T) {
	a := NewAlerter(nil)

	entries := []portscanner.PortEntry{
		makeEntry("0.0.0.0", 8080, 1234, "tcp"),
	}

	alerts := a.Diff(entries)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Type != AlertNewBinding {
		t.Errorf("expected AlertNewBinding, got %s", alerts[0].Type)
	}
}

func TestDiff_NoAlertForKnownEntry(t *testing.T) {
	a := NewAlerter(nil)

	entries := []portscanner.PortEntry{
		makeEntry("0.0.0.0", 8080, 1234, "tcp"),
	}

	// first pass — learns the entry
	a.Diff(entries)
	// second pass — same entry, no new alert expected
	alerts := a.Diff(entries)

	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts for known entry, got %d", len(alerts))
	}
}

func TestDiff_MultipleNewBindings(t *testing.T) {
	a := NewAlerter(nil)

	entries := []portscanner.PortEntry{
		makeEntry("0.0.0.0", 80, 100, "tcp"),
		makeEntry("0.0.0.0", 443, 101, "tcp"),
	}

	alerts := a.Diff(entries)
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestEmit_WritesToOutput(t *testing.T) {
	var buf bytes.Buffer
	a := NewAlerter(&buf)

	entries := []portscanner.PortEntry{
		makeEntry("127.0.0.1", 9090, 555, "tcp"),
	}
	alerts := a.Diff(entries)
	for _, al := range alerts {
		a.Emit(al)
	}

	output := buf.String()
	if !strings.Contains(output, "NEW_BINDING") {
		t.Errorf("expected NEW_BINDING in output, got: %s", output)
	}
	if !strings.Contains(output, "9090") {
		t.Errorf("expected port 9090 in output, got: %s", output)
	}
}
