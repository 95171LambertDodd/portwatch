package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Message holds a single notification payload.
type Message struct {
	Timestamp time.Time
	Level     Level
	Title     string
	Body      string
}

// Notifier dispatches notifications to one or more sinks.
type Notifier struct {
	sinks []Sink
}

// Sink is anything that can receive a Message.
type Sink interface {
	Send(msg Message) error
}

// New creates a Notifier with the provided sinks.
// If no sinks are given, a StdoutSink is added as the default.
func New(sinks ...Sink) *Notifier {
	if len(sinks) == 0 {
		sinks = []Sink{NewStdoutSink(os.Stdout)}
	}
	return &Notifier{sinks: sinks}
}

// Notify sends a message at the given level to all registered sinks.
// Errors from individual sinks are printed to stderr but do not stop delivery.
func (n *Notifier) Notify(level Level, title, body string) {
	msg := Message{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Title:     title,
		Body:      body,
	}
	for _, s := range n.sinks {
		if err := s.Send(msg); err != nil {
			fmt.Fprintf(os.Stderr, "notify: sink error: %v\n", err)
		}
	}
}
