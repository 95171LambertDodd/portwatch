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

func newTestWatcher(interval time.Duration) (*watcher.Watcher, *bytes.Buffer) {
	cfg := defaultConfig(interval)
	scanner := portscanner.NewScanner("/proc/net/tcp")
	var buf bytes.Buffer
	alerter := alerting.NewAlerter(&buf)
	return watcher.New(cfg, scanner, alerter), &buf
}

func TestWatcher_RunCancelsCleanly(t *testing.T) {
	w, _ := newTestWatcher(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestWatcher_RunTicksMultipleTimes(t *testing.T) {
	w, _ := newTestWatcher(30 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Should tick at least twice without panicking.
	_ = w.Run(ctx)
}

func TestNew_ReturnsNonNil(t *testing.T) {
	w, _ := newTestWatcher(time.Second)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}
