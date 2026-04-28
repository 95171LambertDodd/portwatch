// Package projection provides field-level projection for port scan entries,
// allowing callers to select a subset of fields for downstream processing,
// export, or display without modifying the original data.
package projection

import (
	"errors"
	"fmt"

	"github.com/yourorg/portwatch/internal/portscanner"
)

// Field names that can be projected.
const (
	FieldPort     = "port"
	FieldProtocol = "protocol"
	FieldPID      = "pid"
	FieldProcess  = "process"
	FieldCmdline  = "cmdline"
	FieldAddr     = "addr"
)

var validFields = map[string]struct{}{
	FieldPort:     {},
	FieldProtocol: {},
	FieldPID:      {},
	FieldProcess:  {},
	FieldCmdline:  {},
	FieldAddr:     {},
}

// Record is a projected view of a port scan entry containing only the
// requested fields. Absent fields are represented as zero values.
type Record struct {
	Port     uint16 `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	PID      int    `json:"pid,omitempty"`
	Process  string `json:"process,omitempty"`
	Cmdline  string `json:"cmdline,omitempty"`
	Addr     string `json:"addr,omitempty"`
}

// Projector selects a fixed set of fields from port scan entries.
type Projector struct {
	fields map[string]struct{}
}

// New creates a Projector that retains only the specified fields.
// Returns an error if the fields slice is empty or contains an unknown field name.
func New(fields []string) (*Projector, error) {
	if len(fields) == 0 {
		return nil, errors.New("projection: at least one field must be specified")
	}

	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if _, ok := validFields[f]; !ok {
			return nil, fmt.Errorf("projection: unknown field %q", f)
		}
		set[f] = struct{}{}
	}

	return &Projector{fields: set}, nil
}

// Apply returns a Record containing only the fields configured in the Projector.
func (p *Projector) Apply(entry portscanner.PortEntry) Record {
	var r Record

	if _, ok := p.fields[FieldPort]; ok {
		r.Port = entry.Port
	}
	if _, ok := p.fields[FieldProtocol]; ok {
		r.Protocol = entry.Protocol
	}
	if _, ok := p.fields[FieldPID]; ok {
		r.PID = entry.PID
	}
	if _, ok := p.fields[FieldProcess]; ok {
		r.Process = entry.Process
	}
	if _, ok := p.fields[FieldCmdline]; ok {
		r.Cmdline = entry.Cmdline
	}
	if _, ok := p.fields[FieldAddr]; ok {
		r.Addr = entry.Addr
	}

	return r
}

// ApplyAll applies the projection to a slice of entries and returns the
// corresponding slice of Records.
func (p *Projector) ApplyAll(entries []portscanner.PortEntry) []Record {
	out := make([]Record, len(entries))
	for i, e := range entries {
		out[i] = p.Apply(e)
	}
	return out
}

// Fields returns the set of field names this Projector retains.
func (p *Projector) Fields() []string {
	out := make([]string, 0, len(p.fields))
	for f := range p.fields {
		out = append(out, f)
	}
	return out
}
