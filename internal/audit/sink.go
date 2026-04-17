package audit

import "fmt"

// Sink is a destination that can receive audit events.
type Sink interface {
	Receive(Event) error
}

// LoggerSink wraps a Logger to implement the Sink interface.
type LoggerSink struct {
	logger *Logger
}

// NewLoggerSink creates a Sink backed by the given Logger.
func NewLoggerSink(l *Logger) *LoggerSink {
	return &LoggerSink{logger: l}
}

// Receive writes the event to the underlying logger.
func (s *LoggerSink) Receive(e Event) error {
	return s.logger.Log(e)
}

// MultiSink fans out events to multiple sinks.
type MultiSink struct {
	sinks []Sink
}

// NewMultiSink creates a MultiSink from the provided sinks.
func NewMultiSink(sinks ...Sink) *MultiSink {
	return &MultiSink{sinks: sinks}
}

// Receive delivers the event to all sinks, collecting errors.
func (m *MultiSink) Receive(e Event) error {
	var errs []error
	for _, s := range m.sinks {
		if err := s.Receive(e); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("audit: %d sink(s) failed: %v", len(errs), errs[0])
	}
	return nil
}
