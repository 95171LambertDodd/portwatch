// Package export provides functionality to export port scan snapshots
// to various formats (JSON, CSV) for external consumption.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Record is a flattened, exportable representation of a port entry.
type Record struct {
	Timestamp string `json:"timestamp"`
	Proto     string `json:"proto"`
	LocalAddr string `json:"local_addr"`
	Port      uint16 `json:"port"`
	PID       int    `json:"pid"`
	Process   string `json:"process"`
}

// Exporter writes scan results to an output sink.
type Exporter struct {
	format string
}

// New returns an Exporter for the given format ("json" or "csv").
func New(format string) (*Exporter, error) {
	switch format {
	case "json", "csv":
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Write serialises entries to w using the configured format.
func (e *Exporter) Write(w io.Writer, entries []portscanner.PortEntry) error {
	records := toRecords(entries)
	switch e.format {
	case "json":
		return writeJSON(w, records)
	case "csv":
		return writeCSV(w, records)
	}
	return nil
}

func toRecords(entries []portscanner.PortEntry) []Record {
	ts := time.Now().UTC().Format(time.RFC3339)
	out := make([]Record, 0, len(entries))
	for _, e := range entries {
		out = append(out, Record{
			Timestamp: ts,
			Proto:     e.Proto,
			LocalAddr: e.LocalAddr,
			Port:      e.Port,
			PID:       e.PID,
			Process:   e.Process,
		})
	}
	return out
}

func writeJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func writeCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "proto", "local_addr", "port", "pid", "process"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{
			r.Timestamp,
			r.Proto,
			r.LocalAddr,
			fmt.Sprintf("%d", r.Port),
			fmt.Sprintf("%d", r.PID),
			r.Process,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
