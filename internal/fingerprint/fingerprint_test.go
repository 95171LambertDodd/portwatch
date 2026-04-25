package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/portscanner"
)

func makeEntry(proto, addr string, port, pid int) portscanner.PortEntry {
	return portscanner.PortEntry{
		Protocol:  proto,
		LocalAddr: addr,
		LocalPort: port,
		PID:       pid,
	}
}

func TestCompute_ReturnsNonEmpty(t *testing.T) {
	e := makeEntry("tcp", "0.0.0.0", 8080, 1234)
	fp := fingerprint.Compute(e)
	if fp == "" {
		t.Fatal("expected non-empty fingerprint")
	}
}

func TestCompute_Deterministic(t *testing.T) {
	e := makeEntry("tcp", "0.0.0.0", 8080, 1234)
	if fingerprint.Compute(e) != fingerprint.Compute(e) {
		t.Fatal("fingerprint should be deterministic")
	}
}

func TestCompute_DifferentPID_DifferentFingerprint(t *testing.T) {
	a := makeEntry("tcp", "0.0.0.0", 8080, 1000)
	b := makeEntry("tcp", "0.0.0.0", 8080, 2000)
	if fingerprint.Compute(a) == fingerprint.Compute(b) {
		t.Fatal("different PIDs should produce different fingerprints")
	}
}

func TestChanged_FirstCall_ReturnsTrue(t *testing.T) {
	tr := fingerprint.New()
	e := makeEntry("tcp", "0.0.0.0", 9090, 42)
	if !tr.Changed(e) {
		t.Fatal("first observation should always report changed")
	}
}

func TestChanged_SameEntry_ReturnsFalse(t *testing.T) {
	tr := fingerprint.New()
	e := makeEntry("tcp", "127.0.0.1", 3306, 99)
	tr.Changed(e) // first call
	if tr.Changed(e) {
		t.Fatal("identical second call should not report changed")
	}
}

func TestChanged_PIDChange_ReturnsTrue(t *testing.T) {
	tr := fingerprint.New()
	a := makeEntry("tcp", "0.0.0.0", 80, 100)
	b := makeEntry("tcp", "0.0.0.0", 80, 200)
	tr.Changed(a)
	if !tr.Changed(b) {
		t.Fatal("PID change on same port should report changed")
	}
}

func TestRemove_ThenReappear_ReportsChanged(t *testing.T) {
	tr := fingerprint.New()
	e := makeEntry("udp", "0.0.0.0", 53, 7)
	tr.Changed(e)
	tr.Remove(e)
	if !tr.Changed(e) {
		t.Fatal("entry should be treated as new after Remove")
	}
}

func TestLen_TracksCorrectly(t *testing.T) {
	tr := fingerprint.New()
	if tr.Len() != 0 {
		t.Fatalf("expected 0, got %d", tr.Len())
	}
	tr.Changed(makeEntry("tcp", "0.0.0.0", 80, 1))
	tr.Changed(makeEntry("tcp", "0.0.0.0", 443, 2))
	if tr.Len() != 2 {
		t.Fatalf("expected 2, got %d", tr.Len())
	}
	tr.Remove(makeEntry("tcp", "0.0.0.0", 80, 1))
	if tr.Len() != 1 {
		t.Fatalf("expected 1 after remove, got %d", tr.Len())
	}
}
