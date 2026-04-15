package portscanner

import (
	"os"
	"testing"
)

func TestNewScanner(t *testing.T) {
	s := NewScanner()
	if s == nil {
		t.Fatal("expected non-nil scanner")
	}
}

func TestParseProcNet_NonExistentFile(t *testing.T) {
	_, err := parseProcNet("/nonexistent/path/tcp")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestParseProcNet_ValidContent(t *testing.T) {
	// Create a temporary file mimicking /proc/net/tcp format
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 00000000:0050 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 67890 1 0000000000000000 100 0 0 10 0
`
	tmpFile, err := os.CreateTemp("", "proc_net_tcp_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	entries, err := parseProcNet(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only LISTEN (0A) entries should be returned
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.LocalPort != 8080 {
		t.Errorf("expected port 8080 (0x1F90), got %d", entry.LocalPort)
	}
	if entry.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", entry.Protocol)
	}
	if entry.State != "LISTEN" {
		t.Errorf("expected state LISTEN, got %s", entry.State)
	}
}

func TestParseProcNet_EmptyFile(t *testing.T) {
	// A file with only the header line should return zero entries
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
`
	tmpFile, err := os.CreateTemp("", "proc_net_tcp_empty_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	entries, err := parseProcNet(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseProcNet_MultipleListenPorts(t *testing.T) {
	// Both entries are in LISTEN state (0A); both should be returned
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 67890 1 0000000000000000 100 0 0 10 0
`
	tmpFile, err := os.CreateTemp("", "proc_net_tcp_multi_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	entries, err := parseProcNet(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestScan_ReturnsNoError(t *testing.T) {
	s := NewScanner()
	// Scan should not return an error even if /proc/net/tcp doesn't exist
	_, err := s.Scan()
	if err != nil {
		t.Errorf("Scan() returned unexpected error: %v", err)
	}
}
