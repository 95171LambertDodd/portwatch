package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ScanInterval != 5*time.Second {
		t.Errorf("expected 5s scan interval, got %v", cfg.ScanInterval)
	}
	if cfg.AlertOutput != "stdout" {
		t.Errorf("expected stdout alert output, got %q", cfg.AlertOutput)
	}
	if cfg.ProcNetTCP != "/proc/net/tcp" {
		t.Errorf("unexpected ProcNetTCP default: %q", cfg.ProcNetTCP)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	_, err := Load("/tmp/portwatch_nonexistent_12345.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	content := []byte(`
scan_interval: 10s
alert_output: stderr
watched_ports: [80, 443]
ignored_ports: [22]
`)
	f, err := os.CreateTemp("", "portwatch-cfg-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write(content)
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.ScanInterval)
	}
	if len(cfg.WatchedPorts) != 2 {
		t.Errorf("expected 2 watched ports, got %d", len(cfg.WatchedPorts))
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	content := []byte("scan_interval: -1s\n")
	f, err := os.CreateTemp("", "portwatch-cfg-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write(content)
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Fatal("expected error for negative scan_interval")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	content := []byte("scan_interval: [not: valid: yaml\n")
	f, err := os.CreateTemp("", "portwatch-cfg-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write(content)
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestIsIgnored(t *testing.T) {
	cfg := &Config{IgnoredPorts: []int{22, 8080}}
	if !cfg.IsIgnored(22) {
		t.Error("expected port 22 to be ignored")
	}
	if cfg.IsIgnored(80) {
		t.Error("expected port 80 not to be ignored")
	}
}

func TestIsWatched_AllWhenEmpty(t *testing.T) {
	cfg := &Config{WatchedPorts: []int{}}
	if !cfg.IsWatched(9999) {
		t.Error("expected all ports watched when list is empty")
	}
}

func TestIsWatched_Specific(t *testing.T) {
	cfg := &Config{WatchedPorts: []int{80, 443}}
	if !cfg.IsWatched(443) {
		t.Error("expected port 443 to be watched")
	}
	if cfg.IsWatched(8080) {
		t.Error("expected port 8080 not to be watched")
	}
}
