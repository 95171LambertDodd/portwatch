package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portscanner"
)

func sampleEntries() []portscanner.PortEntry {
	return []portscanner.PortEntry{
		{Proto: "tcp", LocalAddr: "0.0.0.0", Port: 8080, PID: 1234, Process: "myapp"},
		{Proto: "udp", LocalAddr: "127.0.0.1", Port: 5353, PID: 5678, Process: "dnsmasq"},
	}
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []string{"json", "csv"} {
		e, err := New(f)
		if err != nil {
			t.Fatalf("New(%q) unexpected error: %v", f, err)
		}
		if e == nil {
			t.Fatalf("New(%q) returned nil", f)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_JSON_ContainsFields(t *testing.T) {
	e, _ := New("json")
	var buf bytes.Buffer
	if err := e.Write(&buf, sampleEntries()); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	var records []Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", records[0].Port)
	}
	if records[1].Process != "dnsmasq" {
		t.Errorf("expected process dnsmasq, got %s", records[1].Process)
	}
}

func TestWrite_CSV_HasHeaderAndRows(t *testing.T) {
	e, _ := New("csv")
	var buf bytes.Buffer
	if err := e.Write(&buf, sampleEntries()); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header+2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("first line should be header, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "8080") {
		t.Errorf("expected port 8080 in row 1: %s", lines[1])
	}
}

func TestWrite_JSON_EmptyEntries(t *testing.T) {
	e, _ := New("json")
	var buf bytes.Buffer
	if err := e.Write(&buf, nil); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got: %s", buf.String())
	}
}
