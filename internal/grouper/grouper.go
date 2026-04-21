// Package grouper aggregates port scan entries by a configurable key
// (e.g. protocol, process name, or port range bucket) so that downstream
// consumers can reason about clusters of activity rather than individual
// bindings.
package grouper

import (
	"fmt"

	"github.com/user/portwatch/internal/portscanner"
)

// GroupBy controls which field is used as the grouping key.
type GroupBy string

const (
	GroupByProtocol GroupBy = "protocol"
	GroupByProcess  GroupBy = "process"
	GroupByPortBand GroupBy = "portband"
)

// Group holds all entries that share the same key.
type Group struct {
	Key     string
	Entries []portscanner.PortEntry
}

// Grouper partitions a slice of PortEntry values into named groups.
type Grouper struct {
	by GroupBy
}

// New returns a Grouper that partitions entries by the given field.
// Returns an error for unrecognised GroupBy values.
func New(by GroupBy) (*Grouper, error) {
	switch by {
	case GroupByProtocol, GroupByProcess, GroupByPortBand:
		return &Grouper{by: by}, nil
	default:
		return nil, fmt.Errorf("grouper: unknown group-by field %q", by)
	}
}

// Group partitions entries and returns a slice of Group values sorted by key.
func (g *Grouper) Group(entries []portscanner.PortEntry) []Group {
	buckets := make(map[string][]portscanner.PortEntry)
	for _, e := range entries {
		k := g.keyFor(e)
		buckets[k] = append(buckets[k], e)
	}

	groups := make([]Group, 0, len(buckets))
	for k, v := range buckets {
		groups = append(groups, Group{Key: k, Entries: v})
	}
	return groups
}

func (g *Grouper) keyFor(e portscanner.PortEntry) string {
	switch g.by {
	case GroupByProtocol:
		return e.Protocol
	case GroupByProcess:
		if e.ProcessName == "" {
			return "unknown"
		}
		return e.ProcessName
	case GroupByPortBand:
		return portBand(e.Port)
	default:
		return "unknown"
	}
}

// portBand returns a human-readable band label for a port number.
func portBand(port int) string {
	switch {
	case port < 1024:
		return "well-known (0-1023)"
	case port < 49152:
		return "registered (1024-49151)"
	default:
		return "dynamic (49152-65535)"
	}
}
