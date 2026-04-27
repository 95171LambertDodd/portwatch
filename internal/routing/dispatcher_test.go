package routing_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/routing"
)

type captureSink struct {
	entries []portscanner.PortEntry
	errOnWrite error
}

func (c *captureSink) Write(e portscanner.PortEntry) error {
	if c.errOnWrite != nil {
		return c.errOnWrite
	}
	c.entries = append(c.entries, e)
	return nil
}

func buildDispatcher(t *testing.T, dead *bytes.Buffer) (*routing.Dispatcher, *captureSink) {
	t.Helper()
	r, err := routing.New([]routing.Rule{
		{Name: "tcp", Protocol: "tcp", Destination: "tcp-sink"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	sink := &captureSink{}
	d, err := routing.NewDispatcher(r, map[string]routing.Writer{"tcp-sink": sink}, dead)
	if err != nil {
		t.Fatalf("NewDispatcher: %v", err)
	}
	return d, sink
}

func TestNewDispatcher_NilRouter_ReturnsError(t *testing.T) {
	_, err := routing.NewDispatcher(nil, map[string]routing.Writer{"s": &captureSink{}}, nil)
	if err == nil {
		t.Fatal("expected error for nil router")
	}
}

func TestNewDispatcher_MissingSink_ReturnsError(t *testing.T) {
	r, _ := routing.New([]routing.Rule{{Name: "r", Destination: "missing"}})
	_, err := routing.NewDispatcher(r, map[string]routing.Writer{"other": &captureSink{}}, nil)
	if err == nil {
		t.Fatal("expected error for missing sink")
	}
}

func TestDispatch_RoutesToCorrectSink(t *testing.T) {
	d, sink := buildDispatcher(t, nil)
	e := entry("tcp", 80, "")
	if err := d.Dispatch(e); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if len(sink.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(sink.entries))
	}
}

func TestDispatch_UnroutedWritesToDeadLetter(t *testing.T) {
	var buf bytes.Buffer
	d, _ := buildDispatcher(t, &buf)
	if err := d.Dispatch(entry("udp", 53, "")); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
	if !strings.Contains(buf.String(), "unrouted") {
		t.Errorf("expected dead-letter log, got %q", buf.String())
	}
}

func TestDispatch_SinkError_Propagated(t *testing.T) {
	d, sink := buildDispatcher(t, nil)
	sink.errOnWrite = errors.New("write failed")
	if err := d.Dispatch(entry("tcp", 443, "")); err == nil {
		t.Fatal("expected sink error to propagate")
	}
}
