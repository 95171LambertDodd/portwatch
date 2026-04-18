// Package pipeline wires scanner output through filter, dedup, suppress, and notify.
package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/dedup"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/suppress"
)

// Pipeline processes port scan results end-to-end.
type Pipeline struct {
	scanner  *portscanner.Scanner
	alerter  *alerting.Alerter
	filter   *filter.Filter
	dedup    *dedup.Deduplicator
	suppress *suppress.Suppressor
	notifier *notify.Notifier
	metrics  *metrics.Metrics
}

// Config holds dependencies for building a Pipeline.
type Config struct {
	Scanner  *portscanner.Scanner
	Alerter  *alerting.Alerter
	Filter   *filter.Filter
	Dedup    *dedup.Deduplicator
	Suppress *suppress.Suppressor
	Notifier *notify.Notifier
	Metrics  *metrics.Metrics
}

// New constructs a Pipeline from the provided Config.
func New(cfg Config) (*Pipeline, error) {
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("pipeline: scanner is required")
	}
	if cfg.Notifier == nil {
		return nil, fmt.Errorf("pipeline: notifier is required")
	}
	return &Pipeline{
		scanner:  cfg.Scanner,
		alerter:  cfg.Alerter,
		filter:   cfg.Filter,
		dedup:    cfg.Dedup,
		suppress: cfg.Suppress,
		notifier: cfg.Notifier,
		metrics:  cfg.Metrics,
	}, nil
}

// Run executes one scan cycle and dispatches alerts.
func (p *Pipeline) Run(ctx context.Context) error {
	entries, err := p.scanner.Scan()
	if err != nil {
		return fmt.Errorf("pipeline: scan failed: %w", err)
	}
	if p.metrics != nil {
		p.metrics.RecordScan()
	}

	var alerts []alerting.Alert
	if p.alerter != nil {
		alerts = p.alerter.Diff(entries)
	}

	for _, a := range alerts {
		if p.filter != nil && !p.filter.Allow(a.Entry) {
			continue
		}
		if p.suppress != nil && p.suppress.IsSuppressed(a.Entry) {
			continue
		}
		if p.dedup != nil && p.dedup.IsDuplicate(a.Entry) {
			continue
		}
		if p.metrics != nil {
			p.metrics.RecordAlert()
		}
		if err := p.notifier.Notify(ctx, a.Entry); err != nil {
			return fmt.Errorf("pipeline: notify failed: %w", err)
		}
	}
	return nil
}
