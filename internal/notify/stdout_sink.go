package notify

import (
	"fmt"
	"io"
)

// StdoutSink writes formatted notifications to an io.Writer.
type StdoutSink struct {
	w io.Writer
}

// NewStdoutSink returns a StdoutSink that writes to w.
func NewStdoutSink(w io.Writer) *StdoutSink {
	return &StdoutSink{w: w}
}

// Send formats msg and writes it to the underlying writer.
func (s *StdoutSink) Send(msg Message) error {
	_, err := fmt.Fprintf(
		s.w,
		"[%s] %s | %s: %s\n",
		msg.Timestamp.Format("2006-01-02T15:04:05Z"),
		msg.Level,
		msg.Title,
		msg.Body,
	)
	return err
}
