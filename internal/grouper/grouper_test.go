package grouper_test

import (
	"sort"
	"testing"

	"github.com/user/portwatch/internal/grouper"
	"github.com/user/portwatch/internal/portscanner"
)

func entry(proto string, port int, process string) portscanner.PortEntry {
	return portscanner.PortEntry{
		Protocol:    proto,
		Port:        port,
		ProcessName: process,
	}
}

func sortGroups(gs []grouper.Group) []grouper.Group {
	sort.Slice(gs, func(i, j int) bool { return gs[i].Key < gs[j].Key })
	return gs
}

func TestNew_ValidGroupBy_ReturnsNonNil(t *testing.T) {
	for _, by := range []grouper.GroupBy{
		grouper.GroupByProtocol,
		grouper.GroupByProcess,
		grouper.GroupByPortBand,
	} {
		g, err := grouper.New(by)
		if err != nil {
			t.Fatalf("New(%q) unexpected error: %v", by, err)
		}
		if g == nil {
			t.Fatalf("New(%q) returned nil", by)
		}
	}
}

func TestNew_InvalidGroupBy_ReturnsError(t *testing.T) {
	_, err := grouper.New("invalid")
	if err == nil {
		t.Fatal("expected error for unknown group-by field, got nil")
	}
}

func TestGroup_ByProtocol(t *testing.T) {
	g, _ := grouper.New(grouper.GroupByProtocol)
	entries := []portscanner.PortEntry{
		entry("tcp", 80, "nginx"),
		entry("udp", 53, "dnsmasq"),
		entry("tcp", 443, "nginx"),
	}
	groups := sortGroups(g.Group(entries))
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Key != "tcp" || len(groups[0].Entries) != 2 {
		t.Errorf("unexpected tcp group: %+v", groups[0])
	}
	if groups[1].Key != "udp" || len(groups[1].Entries) != 1 {
		t.Errorf("unexpected udp group: %+v", groups[1])
	}
}

func TestGroup_ByProcess_UnknownFallback(t *testing.T) {
	g, _ := grouper.New(grouper.GroupByProcess)
	entries := []portscanner.PortEntry{
		entry("tcp", 8080, ""),
		entry("tcp", 9090, ""),
		entry("tcp", 443, "nginx"),
	}
	groups := sortGroups(g.Group(entries))
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	for _, grp := range groups {
		if grp.Key == "unknown" && len(grp.Entries) != 2 {
			t.Errorf("expected 2 unknown-process entries, got %d", len(grp.Entries))
		}
	}
}

func TestGroup_ByPortBand(t *testing.T) {
	g, _ := grouper.New(grouper.GroupByPortBand)
	entries := []portscanner.PortEntry{
		entry("tcp", 80, "nginx"),
		entry("tcp", 8080, "app"),
		entry("tcp", 60000, "ephemeral"),
	}
	groups := g.Group(entries)
	if len(groups) != 3 {
		t.Fatalf("expected 3 port-band groups, got %d", len(groups))
	}
}

func TestGroup_EmptyInput_ReturnsEmptySlice(t *testing.T) {
	g, _ := grouper.New(grouper.GroupByProtocol)
	groups := g.Group(nil)
	if len(groups) != 0 {
		t.Fatalf("expected empty groups, got %d", len(groups))
	}
}
