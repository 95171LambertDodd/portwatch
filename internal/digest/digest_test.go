package digest_test

import (
	"testing"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/portscanner"
)

func entries(ports ...int) []portscanner.Entry {
	out := make([]portscanner.Entry, len(ports))
	for i, p := range ports {
		out[i] = portscanner.Entry{Protocol: "tcp", LocalAddr: "0.0.0.0", LocalPort: p, PID: 0}
	}
	return out
}

func TestCompute_ReturnsNonEmpty(t *testing.T) {
	fp, err := digest.Compute(entries(80, 443))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp == "" {
		t.Fatal("expected non-empty fingerprint")
	}
}

func TestCompute_DeterministicOrderIndependent(t *testing.T) {
	a, err := digest.Compute(entries(80, 443))
	if err != nil {
		t.Fatal(err)
	}
	b, err := digest.Compute(entries(443, 80))
	if err != nil {
		t.Fatal(err)
	}
	if !digest.Equal(a, b) {
		t.Errorf("expected equal fingerprints, got %s vs %s", a, b)
	}
}

func TestCompute_DifferentEntriesProduceDifferentFingerprints(t *testing.T) {
	a, _ := digest.Compute(entries(80))
	b, _ := digest.Compute(entries(8080))
	if digest.Equal(a, b) {
		t.Error("expected different fingerprints for different entries")
	}
}

func TestCompute_EmptySlice(t *testing.T) {
	a, err := digest.Compute([]portscanner.Entry{})
	if err != nil {
		t.Fatal(err)
	}
	b, _ := digest.Compute([]portscanner.Entry{})
	if !digest.Equal(a, b) {
		t.Error("empty slices should produce equal fingerprints")
	}
}

func TestEqual_Reflexive(t *testing.T) {
	fp, _ := digest.Compute(entries(22))
	if !digest.Equal(fp, fp) {
		t.Error("fingerprint should equal itself")
	}
}
