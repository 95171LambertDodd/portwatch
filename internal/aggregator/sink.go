package aggregator

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// SinkRecord is the JSON-serialisable form of a completed Bucket.
type SinkRecord struct {
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Count   int       `json:"count"`
	Ports   []int     `json:"ports"`
}

// FlushSink periodically flushes an Aggregator and writes bucket summaries
// to the provided writer as newline-delimited JSON.
type FlushSink struct {
	agg    *Aggregator
	out    io.Writer
	interval time.Duration
}

// NewFlushSink creates a FlushSink that drains agg every interval.
func NewFlushSink(agg *Aggregator, out io.Writer, interval time.Duration) (*FlushSink, error) {
	if agg == nil {
		return nil, fmt.Errorf("aggregator: FlushSink requires non-nil aggregator")
	}
	if out == nil {
		return nil, fmt.Errorf("aggregator: FlushSink requires non-nil writer")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("aggregator: FlushSink interval must be positive")
	}
	return &FlushSink{agg: agg, out: out, interval: interval}, nil
}

// Flush drains completed buckets and writes each as a JSON line.
func (s *FlushSink) Flush() error {
	buckets := s.agg.Flush()
	for _, b := range buckets {
		rec := toSinkRecord(b)
		data, err := json.Marshal(rec)
		if err != nil {
			return fmt.Errorf("aggregator: marshal bucket: %w", err)
		}
		if _, err := fmt.Fprintf(s.out, "%s\n", data); err != nil {
			return fmt.Errorf("aggregator: write bucket: %w", err)
		}
	}
	return nil
}

func toSinkRecord(b Bucket) SinkRecord {
	ports := make([]int, 0, len(b.Entries))
	seen := map[int]struct{}{}
	for _, e := range b.Entries {
		if _, ok := seen[e.Port]; !ok {
			ports = append(ports, e.Port)
			seen[e.Port] = struct{}{}
		}
	}
	return SinkRecord{
		Start: b.Start,
		End:   b.End,
		Count: b.Count,
		Ports: ports,
	}
}
