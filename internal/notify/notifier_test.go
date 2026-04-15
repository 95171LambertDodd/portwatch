package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/example/portwatch/internal/notify"
)

// errorSink is a Sink that always returns an error, used to test resilience.
type errorSink struct{}

func (e *errorSink) Send(_ notify.Message) error {
	return errors.New("sink failure")
}

func TestNew_DefaultsToStdoutSink(t *testing.T) {
	n := notify.New()
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

func TestNotify_WritesToSink(t *testing.T) {
	var buf bytes.Buffer
	sink := notify.NewStdoutSink(&buf)
	n := notify.New(sink)

	n.Notify(notify.LevelAlert, "port conflict", "port 8080 already bound")

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "port conflict") {
		t.Errorf("expected title in output, got: %s", out)
	}
	if !strings.Contains(out, "port 8080 already bound") {
		t.Errorf("expected body in output, got: %s", out)
	}
}

func TestNotify_MultiSink_AllReceive(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	sink1 := notify.NewStdoutSink(&buf1)
	sink2 := notify.NewStdoutSink(&buf2)
	n := notify.New(sink1, sink2)

	n.Notify(notify.LevelWarn, "new binding", "0.0.0.0:9090")

	for i, buf := range []*bytes.Buffer{&buf1, &buf2} {
		if !strings.Contains(buf.String(), "new binding") {
			t.Errorf("sink %d did not receive message", i+1)
		}
	}
}

func TestNotify_SinkError_DoesNotPanic(t *testing.T) {
	n := notify.New(&errorSink{})
	// Should not panic even when a sink returns an error.
	n.Notify(notify.LevelInfo, "test", "body")
}
