package correlation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/correlation"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write rules file: %v", err)
	}
	return path
}

func TestLoadFile_EmptyPath_ReturnsEmptyCorrelator(t *testing.T) {
	c, err := correlation.LoadFile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil correlator")
	}
}

func TestLoadFile_NonExistentFile_ReturnsError(t *testing.T) {
	_, err := correlation.LoadFile("/no/such/file.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON_ReturnsError(t *testing.T) {
	path := writeRulesFile(t, `{not valid json`)
	_, err := correlation.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadFile_ValidJSON_MatchesRule(t *testing.T) {
	const content = `{
  "rules": [
    {"port": 22, "protocol": "tcp", "name": "ssh", "expected": true}
  ]
}`
	path := writeRulesFile(t, content)
	c, err := correlation.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := entry(22, "tcp")
	svc, ok := c.Lookup(e)
	if !ok {
		t.Fatal("expected match after loading file")
	}
	if svc.Name != "ssh" {
		t.Errorf("expected ssh, got %s", svc.Name)
	}
	if !c.IsExpected(e) {
		t.Error("expected entry to be marked as expected")
	}
}
