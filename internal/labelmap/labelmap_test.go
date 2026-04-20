package labelmap_test

import (
	"testing"

	"github.com/user/portwatch/internal/labelmap"
	"github.com/user/portwatch/internal/portscanner"
)

func TestNew_ValidEntries_ReturnsNonNil(t *testing.T) {
	m, err := labelmap.New([]labelmap.Entry{
		{Port: 80, Protocol: "tcp", Label: "http"},
		{Port: 443, Protocol: "tcp", Label: "https", Description: "TLS"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Map")
	}
}

func TestNew_InvalidProtocol_ReturnsError(t *testing.T) {
	_, err := labelmap.New([]labelmap.Entry{
		{Port: 80, Protocol: "sctp", Label: "bad"},
	})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestNew_EmptyLabel_ReturnsError(t *testing.T) {
	_, err := labelmap.New([]labelmap.Entry{
		{Port: 80, Protocol: "tcp", Label: "   "},
	})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestLookup_KnownPort_ReturnsEntry(t *testing.T) {
	m, _ := labelmap.New([]labelmap.Entry{
		{Port: 22, Protocol: "tcp", Label: "ssh", Description: "Secure Shell"},
	})
	e, ok := m.Lookup(22, "tcp")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Label != "ssh" {
		t.Errorf("expected label 'ssh', got %q", e.Label)
	}
	if e.Description != "Secure Shell" {
		t.Errorf("expected description 'Secure Shell', got %q", e.Description)
	}
}

func TestLookup_UnknownPort_ReturnsFalse(t *testing.T) {
	m, _ := labelmap.New([]labelmap.Entry{
		{Port: 22, Protocol: "tcp", Label: "ssh"},
	})
	_, ok := m.Lookup(9999, "tcp")
	if ok {
		t.Fatal("expected no entry for unknown port")
	}
}

func TestLookup_ProtocolCaseInsensitive(t *testing.T) {
	m, _ := labelmap.New([]labelmap.Entry{
		{Port: 53, Protocol: "UDP", Label: "dns"},
	})
	_, ok := m.Lookup(53, "udp")
	if !ok {
		t.Fatal("expected case-insensitive protocol lookup to succeed")
	}
}

func TestLookupEntry_UsesPortEntry(t *testing.T) {
	m, _ := labelmap.New([]labelmap.Entry{
		{Port: 8080, Protocol: "tcp", Label: "http-alt"},
	})
	pe := portscanner.PortEntry{Port: 8080, Protocol: "tcp"}
	e, ok := m.LookupEntry(pe)
	if !ok {
		t.Fatal("expected entry via LookupEntry")
	}
	if e.Label != "http-alt" {
		t.Errorf("expected 'http-alt', got %q", e.Label)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	input := []labelmap.Entry{
		{Port: 80, Protocol: "tcp", Label: "http"},
		{Port: 443, Protocol: "tcp", Label: "https"},
		{Port: 53, Protocol: "udp", Label: "dns"},
	}
	m, _ := labelmap.New(input)
	all := m.All()
	if len(all) != len(input) {
		t.Errorf("expected %d entries, got %d", len(input), len(all))
	}
}
