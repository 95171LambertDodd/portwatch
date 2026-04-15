package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Info holds process metadata for a given PID.
type Info struct {
	PID     int
	Name    string
	Cmdline string
}

// Resolver looks up process information by PID from the /proc filesystem.
type Resolver struct {
	procRoot string
}

// NewResolver creates a Resolver using the given proc root (typically "/proc").
func NewResolver(procRoot string) *Resolver {
	if procRoot == "" {
		procRoot = "/proc"
	}
	return &Resolver{procRoot: procRoot}
}

// Lookup returns process Info for the given PID.
// Returns an error if the process cannot be found or read.
func (r *Resolver) Lookup(pid int) (*Info, error) {
	name, err := r.readComm(pid)
	if err != nil {
		return nil, fmt.Errorf("resolver: read comm for pid %d: %w", pid, err)
	}
	cmdline, err := r.readCmdline(pid)
	if err != nil {
		// cmdline is best-effort; don't fail hard
		cmdline = ""
	}
	return &Info{
		PID:     pid,
		Name:    name,
		Cmdline: cmdline,
	}, nil
}

func (r *Resolver) readComm(pid int) (string, error) {
	path := filepath.Join(r.procRoot, strconv.Itoa(pid), "comm")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (r *Resolver) readCmdline(pid int) (string, error) {
	path := filepath.Join(r.procRoot, strconv.Itoa(pid), "cmdline")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	// cmdline is null-byte separated
	return strings.ReplaceAll(strings.TrimRight(string(data), "\x00"), "\x00", " "), nil
}
