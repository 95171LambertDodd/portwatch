package correlation_test

import (
	"testing"

	"github.com/user/portwatch/internal/correlation"
	"github.com/user/portwatch/internal/portscanner"
)

func entry(port uint16, proto string) portscanner.PortEntry {
	return portscanner.PortEntry{Port: port, Protocol: proto}
}

func TestNew_ValidRules_ReturnsNonNil(t *testing.T) {
	c, err := correlation.New([]correlation.Rule{
		{Port: 22, Protocol: "tcp", Service: correlation.ServiceInfo{Name: "ssh"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil correlator")
	}
}

func TestNew_InvalidProtocol_ReturnsError(t *testing.T) {
	_, err := correlation.New([]correlation.Rule{
		{Port: 80, Protocol: "ftp", Service: correlation.ServiceInfo{Name: "bad"}},
	})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestLookup_KnownPort_ReturnsService(t *testing.T) {
	c, _ := correlation.New([]correlation.Rule{
		{Port: 443, Protocol: "tcp", Service: correlation.ServiceInfo{Name: "https", Expected: true}},
	})
	svc, ok := c.Lookup(entry(443, "tcp"))
	if !ok {
		t.Fatal("expected match")
	}
	if svc.Name != "https" {
		t.Errorf("expected https, got %s", svc.Name)
	}
}

func TestLookup_UnknownPort_ReturnsFalse(t *testing.T) {
	c, _ := correlation.New(nil)
	_, ok := c.Lookup(entry(9999, "tcp"))
	if ok {
		t.Fatal("expected no match for unknown port")
	}
}

func TestIsExpected_ExpectedEntry_ReturnsTrue(t *testing.T) {
	c, _ := correlation.New([]correlation.Rule{
		{Port: 22, Protocol: "tcp", Service: correlation.ServiceInfo{Name: "ssh", Expected: true}},
	})
	if !c.IsExpected(entry(22, "tcp")) {
		t.Fatal("expected true for known expected entry")
	}
}

func TestIsExpected_UnexpectedEntry_ReturnsFalse(t *testing.T) {
	c, _ := correlation.New([]correlation.Rule{
		{Port: 8080, Protocol: "tcp", Service: correlation.ServiceInfo{Name: "dev-server", Expected: false}},
	})
	if c.IsExpected(entry(8080, "tcp")) {
		t.Fatal("expected false for not-expected entry")
	}
}

func TestLookup_ProtocolCaseInsensitive(t *testing.T) {
	c, _ := correlation.New([]correlation.Rule{
		{Port: 53, Protocol: "UDP", Service: correlation.ServiceInfo{Name: "dns"}},
	})
	_, ok := c.Lookup(entry(53, "udp"))
	if !ok {
		t.Fatal("expected match with lowercase protocol lookup")
	}
}
