// Package topology provides a Builder that groups observed port entries by
// their owning process (PID), producing a Map of Node values. Each Node
// captures the process identity (PID, comm, cmdline) alongside every port
// it holds, giving operators a clear picture of which processes are
// responsible for which network bindings at any given scan interval.
//
// Usage:
//
//	b := topology.New()
//	m, err := b.Build(entries)
//	if err != nil { ... }
//	for _, node := range m.Nodes {
//	    fmt.Printf("pid=%d comm=%s ports=%d\n", node.PID, node.Comm, len(node.Ports))
//	}
package topology
