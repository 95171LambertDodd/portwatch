package aggregator

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewFlushSink_NilAggregator_ReturnsError(t *testing.T) {
	_, err := NewFlushSink(nil, &bytes.Buffer{}, time.Second)
	if err == nil {
		t.Fatal("expected error for nil aggregator")
	}
}

func TestNewFlushSink_NilWriter_ReturnsError(t *testing.T) {
	a, _ := New(time.Minute, 5)
	_, err := NewFlushSink(a, nil, time.Second)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestNewFlushSink_ZeroInterval_ReturnsError(t *testing.T) {
	a, _ := New(time.Minute, 5)
	_, err := NewFlushSink(a, &bytes.Buffer{}, 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestFlushSink_Flush_WritesJSON(t *testing.T) {
	now := t0
	clock := func() time.Time { return now }
	agg, _ := newWithClock(time.Minute, 5, clock)

	agg.Add(makeEntries(80, 443))
	now = t0.Add(2 * time.Minute)
	agg.Add(makeEntries(22)) // triggers rotation of first bucket

	var buf bytes.Buffer
	sink, err := NewFlushSink(agg, &buf, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := sink.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 JSON line, got %d", len(lines))
	}

	var rec SinkRecord
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if rec.Count != 2 {
		t.Errorf("expected count 2, got %d", rec.Count)
	}
	if len(rec.Ports) != 2 {
		t.Errorf("expected 2 unique ports, got %d", len(rec.Ports))
	}
}

func TestFlushSink_Flush_EmptyAggregator_WritesNothing(t *testing.T) {
	a, _ := New(time.Minute, 5)
	var buf bytes.Buffer
	sink, _ := NewFlushSink(a, &buf, time.Second)
	if err := sink.Flush(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %s", buf.String())
	}
}
