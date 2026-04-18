package rollup

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/notify"
)

// FlushSink periodically flushes the Aggregator and emits summaries via a Notifier.
type FlushSink struct {
	agg      *Aggregator
	notifier *notify.Notifier
	interval time.Duration
}

// NewFlushSink returns a FlushSink that drains agg every interval.
func NewFlushSink(agg *Aggregator, n *notify.Notifier, interval time.Duration) *FlushSink {
	return &FlushSink{agg: agg, notifier: n, interval: interval}
}

// Run starts the flush loop; it blocks until ctx is cancelled.
func (f *FlushSink) Run(ctx context.Context) {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			f.flush()
			return
		case <-ticker.C:
			f.flush()
		}
	}
}

func (f *FlushSink) flush() {
	events := f.agg.Flush()
	for _, e := range events {
		f.notifier.Notify(Summary(e)) //nolint:errcheck
	}
}
