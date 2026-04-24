package topology

import (
	"testing"

	"github.com/example/portwatch/internal/portscanner"
)

func entry(pid int, comm, proto string, port uint16) portscanner.PortEntry {
	return portscanner.PortEntry{
		PID:      pid,
		Comm:     comm,
		Protocol: proto,
		Port:     port,
	}
}

func TestNew_ReturnsNonNil(t *testing.T) {
	b := New()
	if b == nil {
		t.Fatal("expected non-nil Builder")
	}
}

func TestBuild_NilEntries_ReturnsError(t *testing.T) {
	b := New()
	_, err := b.Build(nil)
	if err == nil {
		t.Fatal("expected error for nil entries")
	}
}

func TestBuild_EmptySlice_ReturnsEmptyMap(t *testing.T) {
	b := New()
	m, err := b.Build([]portscanner.PortEntry{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Nodes) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(m.Nodes))
	}
}

func TestBuild_SingleEntry_SingleNode(t *testing.T) {
	b := New()
	entries := []portscanner.PortEntry{entry(42, "nginx", "tcp", 80)}
	m, err := b.Build(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(m.Nodes))
	}
	if m.Nodes[0].PID != 42 {
		t.Errorf("expected PID 42, got %d", m.Nodes[0].PID)
	}
	if m.Nodes[0].Comm != "nginx" {
		t.Errorf("expected comm nginx, got %s", m.Nodes[0].Comm)
	}
	if len(m.Nodes[0].Ports) != 1 {
		t.Errorf("expected 1 port, got %d", len(m.Nodes[0].Ports))
	}
}

func TestBuild_MultipleEntriesSamePID_GroupedUnderOneNode(t *testing.T) {
	b := New()
	entries := []portscanner.PortEntry{
		entry(10, "sshd", "tcp", 22),
		entry(10, "sshd", "tcp", 2222),
	}
	m, err := b.Build(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(m.Nodes))
	}
	if len(m.Nodes[0].Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(m.Nodes[0].Ports))
	}
}

func TestBuild_DifferentPIDs_SeparateNodes_SortedByPID(t *testing.T) {
	b := New()
	entries := []portscanner.PortEntry{
		entry(200, "postgres", "tcp", 5432),
		entry(100, "redis", "tcp", 6379),
	}
	m, err := b.Build(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(m.Nodes))
	}
	if m.Nodes[0].PID != 100 || m.Nodes[1].PID != 200 {
		t.Errorf("nodes not sorted by PID: %d, %d", m.Nodes[0].PID, m.Nodes[1].PID)
	}
}

func TestBuild_ZeroPID_SyntheticCommAssigned(t *testing.T) {
	b := New()
	e := portscanner.PortEntry{PID: 0, Protocol: "tcp", Port: 111}
	m, err := b.Build([]portscanner.PortEntry{e})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Nodes) != 1 {
		t.Fatalf("expected 1 node")
	}
	if m.Nodes[0].Comm == "" {
		t.Error("expected synthetic comm for PID 0")
	}
}
