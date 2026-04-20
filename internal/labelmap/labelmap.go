// Package labelmap provides a registry for mapping port/protocol pairs to
// human-readable service labels used in alert output and history records.
package labelmap

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/portscanner"
)

// Entry associates a port and protocol with a label and optional description.
type Entry struct {
	Port        uint16 `json:"port"`
	Protocol    string `json:"protocol"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// Map holds the registered label entries keyed by port+protocol.
type Map struct {
	entries map[string]Entry
}

// New creates a Map from the given entries. Returns an error if any entry has
// an invalid protocol or empty label.
func New(entries []Entry) (*Map, error) {
	m := &Map{entries: make(map[string]Entry, len(entries))}
	for _, e := range entries {
		proto := strings.ToLower(e.Protocol)
		if proto != "tcp" && proto != "udp" {
			return nil, fmt.Errorf("labelmap: invalid protocol %q for port %d", e.Protocol, e.Port)
		}
		if strings.TrimSpace(e.Label) == "" {
			return nil, fmt.Errorf("labelmap: empty label for port %d/%s", e.Port, proto)
		}
		key := entryKey(e.Port, proto)
		m.entries[key] = Entry{
			Port:        e.Port,
			Protocol:    proto,
			Label:       e.Label,
			Description: e.Description,
		}
	}
	return m, nil
}

// Lookup returns the label entry for the given port/protocol pair.
// The second return value is false when no mapping exists.
func (m *Map) Lookup(port uint16, protocol string) (Entry, bool) {
	e, ok := m.entries[entryKey(port, strings.ToLower(protocol))]
	return e, ok
}

// LookupEntry is a convenience wrapper that accepts a portscanner.PortEntry.
func (m *Map) LookupEntry(pe portscanner.PortEntry) (Entry, bool) {
	return m.Lookup(pe.Port, pe.Protocol)
}

// All returns a slice of all registered entries in undefined order.
func (m *Map) All() []Entry {
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, e)
	}
	return out
}

func entryKey(port uint16, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}
