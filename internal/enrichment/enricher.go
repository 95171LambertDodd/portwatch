// Package enrichment provides port entry enrichment by attaching process,
// correlation, and tag metadata to raw scanner entries before they are
// dispatched to alerting or notification sinks.
package enrichment

import (
	"fmt"

	"github.com/yourorg/portwatch/internal/correlation"
	"github.com/yourorg/portwatch/internal/portscanner"
	"github.com/yourorg/portwatch/internal/process"
	"github.com/yourorg/portwatch/internal/tag"
)

// EnrichedEntry wraps a raw scanner entry with additional metadata resolved
// at enrichment time.
type EnrichedEntry struct {
	portscanner.Entry

	// PID is the process ID bound to the port, or 0 if unresolved.
	PID int
	// ProcessName is the executable name for the owning process.
	ProcessName string
	// Cmdline is the full command line of the owning process.
	Cmdline string
	// Service is the correlated service name, empty if unknown.
	Service string
	// Tags are labels applied by configured tag rules.
	Tags []string
}

// Enricher attaches metadata to raw port entries.
type Enricher struct {
	resolver    *process.Resolver
	correlator  *correlation.Correlator
	tagger      *tag.Tagger
}

// Config holds optional dependencies for the Enricher. Any nil field is
// silently skipped during enrichment.
type Config struct {
	Resolver   *process.Resolver
	Correlator *correlation.Correlator
	Tagger     *tag.Tagger
}

// New creates an Enricher from the provided Config. All fields are optional;
// passing a zero-value Config produces an enricher that copies entries
// without modification.
func New(cfg Config) *Enricher {
	return &Enricher{
		resolver:   cfg.Resolver,
		correlator: cfg.Correlator,
		tagger:     cfg.Tagger,
	}
}

// Enrich resolves process info, service correlation, and tags for a single
// port entry, returning an EnrichedEntry. Errors from individual resolvers
// are non-fatal; the entry is returned with whatever information was
// successfully gathered.
func (e *Enricher) Enrich(entry portscanner.Entry) EnrichedEntry {
	out := EnrichedEntry{Entry: entry}

	if e.resolver != nil && entry.PID > 0 {
		if info, err := e.resolver.Lookup(entry.PID); err == nil {
			out.PID = entry.PID
			out.ProcessName = info.Name
			out.Cmdline = info.Cmdline
		}
	}

	if e.correlator != nil {
		if svc, ok := e.correlator.Lookup(entry.Port, entry.Protocol); ok {
			out.Service = svc
		}
	}

	if e.tagger != nil {
		out.Tags = e.tagger.Tag(entry)
	}

	return out
}

// EnrichAll enriches a slice of entries, returning a corresponding slice of
// EnrichedEntry values in the same order.
func (e *Enricher) EnrichAll(entries []portscanner.Entry) []EnrichedEntry {
	out := make([]EnrichedEntry, 0, len(entries))
	for _, en := range entries {
		out = append(out, e.Enrich(en))
	}
	return out
}

// String returns a human-readable summary of an EnrichedEntry.
func (ee EnrichedEntry) String() string {
	base := fmt.Sprintf("%s:%d", ee.Protocol, ee.Port)
	if ee.Service != "" {
		base += fmt.Sprintf(" [%s]", ee.Service)
	}
	if ee.ProcessName != "" {
		base += fmt.Sprintf(" pid=%d(%s)", ee.PID, ee.ProcessName)
	}
	if len(ee.Tags) > 0 {
		base += fmt.Sprintf(" tags=%v", ee.Tags)
	}
	return base
}
