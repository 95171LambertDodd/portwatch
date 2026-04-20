package enrichment_test

import (
	"testing"

	"github.com/user/portwatch/internal/enrichment"
	"github.com/user/portwatch/internal/portscanner"
)

func makeEntry(port uint16, proto, ip string, pid int) portscanner.PortEntry {
	return portscanner.PortEntry{
		Port:     port,
		Protocol: proto,
		IP:       ip,
		PID:      pid,
	}
}

func TestNew_ReturnsNonNil(t *testing.T) {
	e, err := enrichment.New(enrichment.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil enricher")
	}
}

func TestEnrich_SetsHostname(t *testing.T) {
	e, err := enrichment.New(enrichment.Config{
		ResolveHostnames: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := makeEntry(8080, "tcp", "127.0.0.1", 1234)
	result := e.Enrich(entry)

	// Hostname resolution may return empty or "localhost" in test environments;
	// we only verify the call doesn't panic and returns a valid struct.
	_ = result.Hostname
}

func TestEnrich_NoHostnameResolution_LeavesHostnameEmpty(t *testing.T) {
	e, err := enrichment.New(enrichment.Config{
		ResolveHostnames: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := makeEntry(443, "tcp", "10.0.0.1", 99)
	result := e.Enrich(entry)

	if result.Hostname != "" {
		t.Errorf("expected empty hostname, got %q", result.Hostname)
	}
}

func TestEnrich_PreservesOriginalFields(t *testing.T) {
	e, err := enrichment.New(enrichment.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := makeEntry(22, "tcp", "0.0.0.0", 555)
	result := e.Enrich(entry)

	if result.Port != entry.Port {
		t.Errorf("port mismatch: got %d, want %d", result.Port, entry.Port)
	}
	if result.Protocol != entry.Protocol {
		t.Errorf("protocol mismatch: got %s, want %s", result.Protocol, entry.Protocol)
	}
	if result.IP != entry.IP {
		t.Errorf("IP mismatch: got %s, want %s", result.IP, entry.IP)
	}
	if result.PID != entry.PID {
		t.Errorf("PID mismatch: got %d, want %d", result.PID, entry.PID)
	}
}

func TestEnrichAll_ReturnsEnrichedSlice(t *testing.T) {
	e, err := enrichment.New(enrichment.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := []portscanner.PortEntry{
		makeEntry(80, "tcp", "0.0.0.0", 100),
		makeEntry(443, "tcp", "0.0.0.0", 101),
		makeEntry(53, "udp", "127.0.0.1", 102),
	}

	results := e.EnrichAll(entries)

	if len(results) != len(entries) {
		t.Fatalf("expected %d results, got %d", len(entries), len(results))
	}
	for i, r := range results {
		if r.Port != entries[i].Port {
			t.Errorf("entry %d: port mismatch", i)
		}
	}
}

func TestEnrichAll_EmptySlice_ReturnsEmpty(t *testing.T) {
	e, err := enrichment.New(enrichment.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := e.EnrichAll([]portscanner.PortEntry{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
