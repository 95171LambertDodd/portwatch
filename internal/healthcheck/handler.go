package healthcheck

import (
	"encoding/json"
	"io"
	"time"
)

// Reporter can write a health report to any writer.
type Reporter struct {
	monitor *Monitor
	out     io.Writer
	now     func() time.Time
}

// NewReporter creates a Reporter backed by the given Monitor.
func NewReporter(m *Monitor, out io.Writer) *Reporter {
	return &Reporter{monitor: m, out: out, now: time.Now}
}

// Write runs all checks and serialises the report as JSON to the output writer.
func (r *Reporter) Write() error {
	report := r.monitor.Check()
	report.GeneratedAt = r.now()
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
