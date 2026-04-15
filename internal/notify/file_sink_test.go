package notify_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/notify"
)

func TestFileSink_CreatesFileAndWritesJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")

	sink := notify.NewFileSink(path)
	msg := notify.Message{
		Timestamp: time.Now().UTC(),
		Level:     notify.LevelAlert,
		Title:     "unexpected binding",
		Body:      "127.0.0.1:3306",
	}

	if err := sink.Send(msg); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}

	var got notify.Message
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("could not unmarshal log entry: %v", err)
	}

	if got.Level != notify.LevelAlert {
		t.Errorf("expected level ALERT, got %s", got.Level)
	}
	if got.Title != msg.Title {
		t.Errorf("expected title %q, got %q", msg.Title, got.Title)
	}
}

func TestFileSink_AppendsMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")
	sink := notify.NewFileSink(path)

	for i := 0; i < 3; i++ {
		msg := notify.Message{
			Timestamp: time.Now().UTC(),
			Level:     notify.LevelInfo,
			Title:     "tick",
			Body:      "ok",
		}
		if err := sink.Send(msg); err != nil {
			t.Fatalf("Send[%d] error: %v", i, err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	// Each entry is newline-delimited JSON; expect 3 lines.
	lines := 0
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}
	if lines != 3 {
		t.Errorf("expected 3 log lines, got %d", lines)
	}
}
