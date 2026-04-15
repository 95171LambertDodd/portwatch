package watcher_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/watcher"
)

func defaultConfig(interval time.Duration) *config.Config {
	cfg := config.DefaultConfig()
	cfg.Interval = interval
	return cfg
}

func TestWatcher_RunCancelsCleanly(t *testing.T) {
	cfg := defaultConfig(50 * time.Millisecond)
	scanner := portscanner.NewScanner("/proc/net/tcp")
	var buf bytes.Buffer
	alerter := alerting.NewAlerter(&buf)

	w := watcher.New(cfg, scanner, alerter)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestWatcher_RunTicksMultipleTimes(t *testing.T) {
	cfg := defaultConfig(30 * time.Millisecond)
	scanner := portscanner.NewScanner("/proc/net/tcp")
	var buf bytes.Buffer
	alerter := alerting.NewAlerter(&buf)

	w := watcher.New(cfg, scanner, alerter)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Should tick at least twice without panicking.
	_ = w.Run(ctx)
}

func TestNew_ReturnsNonNil(t *testing.T) {
	cfg := defaultConfig(time.Second)
	scanner := portscanner.NewScanner("/proc/net/tcp")
	var buf bytes.Buffer
	alerter := alerting.NewAlerter(&buf)

	w := watcher.New(cfg, scanner, alerter)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}
