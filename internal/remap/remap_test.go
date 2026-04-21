package remap_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/remap"
)

func entry(port int, proto string) portscanner.PortEntry {
	return portscanner.PortEntry{Port: port, Protocol: proto}
}

func TestNew_EmptyRules_ReturnsNonNil(t *testing.T) {
	r, err := remap.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Remapper")
	}
}

func TestNew_InvalidProtocol_ReturnsError(t *testing.T) {
	_, err := remap.New([]remap.Rule{
		{FromPort: 80, Protocol: "icmp", Alias: "http"},
	})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestNew_EmptyAlias_ReturnsError(t *testing.T) {
	_, err := remap.New([]remap.Rule{
		{FromPort: 80, Protocol: "tcp", Alias: ""},
	})
	if err == nil {
		t.Fatal("expected error for empty alias")
	}
}

func TestLookup_KnownPort_ReturnsMapped(t *testing.T) {
	r, err := remap.New([]remap.Rule{
		{FromPort: 8080, Protocol: "tcp", Alias: "http-alt", ToPort: 80},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := r.Lookup(entry(8080, "tcp"))
	if !res.Mapped {
		t.Fatal("expected Mapped=true")
	}
	if res.Alias != "http-alt" {
		t.Errorf("expected alias http-alt, got %q", res.Alias)
	}
	if res.ToPort != 80 {
		t.Errorf("expected ToPort 80, got %d", res.ToPort)
	}
}

func TestLookup_UnknownPort_ReturnsNotMapped(t *testing.T) {
	r, _ := remap.New([]remap.Rule{
		{FromPort: 443, Protocol: "tcp", Alias: "https"},
	})
	res := r.Lookup(entry(9999, "tcp"))
	if res.Mapped {
		t.Fatal("expected Mapped=false for unknown port")
	}
}

func TestLookup_ProtocolMismatch_ReturnsNotMapped(t *testing.T) {
	r, _ := remap.New([]remap.Rule{
		{FromPort: 53, Protocol: "tcp", Alias: "dns-tcp"},
	})
	res := r.Lookup(entry(53, "udp"))
	if res.Mapped {
		t.Fatal("expected Mapped=false for protocol mismatch")
	}
}

func TestApply_RemapsPort(t *testing.T) {
	r, _ := remap.New([]remap.Rule{
		{FromPort: 8080, Protocol: "tcp", Alias: "http-alt", ToPort: 80},
	})
	out := r.Apply(entry(8080, "tcp"))
	if out.Port != 80 {
		t.Errorf("expected port 80 after apply, got %d", out.Port)
	}
}

func TestApply_NoToPort_LeavesPortUnchanged(t *testing.T) {
	r, _ := remap.New([]remap.Rule{
		{FromPort: 8080, Protocol: "tcp", Alias: "http-alt"},
	})
	out := r.Apply(entry(8080, "tcp"))
	if out.Port != 8080 {
		t.Errorf("expected port 8080 unchanged, got %d", out.Port)
	}
}
