package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/portscanner"
)

func TestNew_MissingScanner_ReturnsError(t *testing.T) {
	n := notify.New(notify.Config{})
	_, err := pipeline.New(pipeline.Config{Notifier: n})
	if err == nil {
		t.Fatal("expected error for missing scanner")
	}
}

func TestNew_MissingNotifier_ReturnsError(t *testing.T) {
	s := portscanner.NewScanner([]string{})
	_, err := pipeline.New(pipeline.Config{Scanner: s})
	if err == nil {
		t.Fatal("expected error for missing notifier")
	}
}

func TestNew_ValidConfig_ReturnsNonNil(t *testing.T) {
	s := portscanner.NewScanner([]string{})
	n := notify.New(notify.Config{})
	p, err := pipeline.New(pipeline.Config{Scanner: s, Notifier: n})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestRun_NoEntries_NoError(t *testing.T) {
	s := portscanner.NewScanner([]string{})
	n := notify.New(notify.Config{})
	p, _ := pipeline.New(pipeline.Config{Scanner: s, Notifier: n})
	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
