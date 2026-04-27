package shadow

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portscanner"
)

func buildSink(t *testing.T, known map[uint16]bool, jsonOut *bytes.Buffer) *Sink {
	t.Helper()
	ev := newWithClock(&bytes.Buffer{}, fixedClock)
	s, err := NewSink(ev, known, jsonOut)
	if err != nil {
		t.Fatalf("NewSink: %v", err)
	}
	return s
}

func TestNewSink_NilEvaluator_ReturnsError(t *testing.T) {
	_, err := NewSink(nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for nil evaluator")
	}
}

func TestNewSink_ValidParams_ReturnsNonNil(t *testing.T) {
	s := buildSink(t, nil, nil)
	if s == nil {
		t.Fatal("expected non-nil sink")
	}
}

func TestSink_Write_UnknownPort_RecordsEvent(t *testing.T) {
	s := buildSink(t, map[uint16]bool{}, nil)
	err := s.Write(portscanner.PortEntry{Port: 4444, Protocol: "tcp", Process: "svc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.eval.Events()) != 1 {
		t.Fatalf("expected 1 event, got %d", len(s.eval.Events()))
	}
}

func TestSink_Write_KnownPort_NoEvent(t *testing.T) {
	s := buildSink(t, map[uint16]bool{443: true}, nil)
	s.Write(portscanner.PortEntry{Port: 443, Protocol: "tcp", Process: "nginx"}) //nolint:errcheck
	if len(s.eval.Events()) != 0 {
		t.Fatalf("expected 0 events for known port, got %d", len(s.eval.Events()))
	}
}

func TestSink_Flush_NilWriter_NoOp(t *testing.T) {
	s := buildSink(t, nil, nil)
	if err := s.Flush(); err != nil {
		t.Fatalf("unexpected error on nil-writer flush: %v", err)
	}
}

func TestSink_Flush_WritesJSON(t *testing.T) {
	var jsonBuf bytes.Buffer
	s := buildSink(t, map[uint16]bool{}, &jsonBuf)
	s.Write(portscanner.PortEntry{Port: 8888, Protocol: "tcp", Process: "app"}) //nolint:errcheck
	if err := s.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBuf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON from Flush: %v", err)
	}
	if !strings.Contains(jsonBuf.String(), "event_count") {
		t.Errorf("expected event_count in JSON output")
	}
	count, ok := result["event_count"].(float64)
	if !ok || int(count) != 1 {
		t.Errorf("expected event_count=1, got %v", result["event_count"])
	}
}
