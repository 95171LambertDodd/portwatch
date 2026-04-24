// Package topology builds a port-to-process dependency map, grouping
// observed port entries by their owning process to reveal which processes
// hold which ports at a given point in time.
package topology

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/example/portwatch/internal/portscanner"
)

// Node represents a single process and all ports it currently holds.
type Node struct {
	PID      int
	Comm     string
	Cmdline  string
	Ports    []portscanner.PortEntry
}

// Map is a snapshot of the port-to-process topology.
type Map struct {
	Nodes []*Node
}

// Builder constructs topology maps from port scan results.
type Builder struct {
	mu sync.Mutex
}

// New returns a new Builder.
func New() *Builder {
	return &Builder{}
}

// Build groups the provided port entries by PID and returns a topology Map.
// Entries with PID == 0 are collected under a synthetic "kernel" node.
func (b *Builder) Build(entries []portscanner.PortEntry) (*Map, error) {
	if entries == nil {
		return nil, errors.New("topology: entries must not be nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	index := make(map[int]*Node)

	for _, e := range entries {
		pid := e.PID
		node, ok := index[pid]
		if !ok {
			comm := e.Comm
			if comm == "" {
				comm = fmt.Sprintf("pid-%d", pid)
			}
			node = &Node{
				PID:     pid,
				Comm:    comm,
				Cmdline: e.Cmdline,
			}
			index[pid] = node
		}
		node.Ports = append(node.Ports, e)
	}

	nodes := make([]*Node, 0, len(index))
	for _, n := range index {
		nodes = append(nodes, n)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].PID < nodes[j].PID
	})

	return &Map{Nodes: nodes}, nil
}
