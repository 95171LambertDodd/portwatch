package process_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/yourorg/portwatch/internal/process"
)

func makeFakeProc(t *testing.T, pid int, comm, cmdline string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, strconv.Itoa(pid))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "comm"), []byte(comm+"\n"), 0644); err != nil {
		t.Fatalf("write comm: %v", err)
	}
	if cmdline != "" {
		nullSep := filepath.Join(dir, "cmdline")
		if err := os.WriteFile(nullSep, []byte(cmdline), 0644); err != nil {
			t.Fatalf("write cmdline: %v", err)
		}
	}
	return root
}

func TestNewResolver_ReturnsNonNil(t *testing.T) {
	r := process.NewResolver("")
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestLookup_ValidPID(t *testing.T) {
	root := makeFakeProc(t, 1234, "nginx", "/usr/sbin/nginx\x00-g\x00daemon off;\x00")
	r := process.NewResolver(root)
	info, err := r.Lookup(1234)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.PID != 1234 {
		t.Errorf("expected PID 1234, got %d", info.PID)
	}
	if info.Name != "nginx" {
		t.Errorf("expected name 'nginx', got %q", info.Name)
	}
	if info.Cmdline != "/usr/sbin/nginx -g daemon off;" {
		t.Errorf("unexpected cmdline: %q", info.Cmdline)
	}
}

func TestLookup_NonExistentPID(t *testing.T) {
	r := process.NewResolver(t.TempDir())
	_, err := r.Lookup(9999)
	if err == nil {
		t.Fatal("expected error for non-existent PID")
	}
}

func TestLookup_MissingCmdline_StillSucceeds(t *testing.T) {
	root := makeFakeProc(t, 42, "sshd", "")
	r := process.NewResolver(root)
	info, err := r.Lookup(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "sshd" {
		t.Errorf("expected 'sshd', got %q", info.Name)
	}
	if info.Cmdline != "" {
		t.Errorf("expected empty cmdline, got %q", info.Cmdline)
	}
}
