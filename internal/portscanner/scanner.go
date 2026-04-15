package portscanner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single bound port entry on the system.
type PortEntry struct {
	Protocol string
	LocalAddress string
	LocalPort int
	PID int
	State string
}

// Scanner reads active port bindings from the system.
type Scanner struct{}

// NewScanner creates a new Scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Scan reads /proc/net/tcp and /proc/net/tcp6 and returns active port entries.
func (s *Scanner) Scan() ([]PortEntry, error) {
	var entries []PortEntry

	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		results, err := parseProcNet(path)
		if err != nil {
			// Non-fatal: file may not exist on all systems
			continue
		}
		entries = append(entries, results...)
	}

	return entries, nil
}

// parseProcNet parses a /proc/net/tcp or /proc/net/tcp6 file.
func parseProcNet(path string) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var entries []PortEntry
	scanner := bufio.NewScanner(f)

	// Skip header line
	if scanner.Scan() {
		_ = scanner.Text()
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		localAddrPort := fields[1]
		state := fields[3]
		// Only LISTEN state (0A)
		if state != "0A" {
			continue
		}

		parts := strings.Split(localAddrPort, ":")
		if len(parts) != 2 {
			continue
		}

		portHex := parts[1]
		port, err := strconv.ParseInt(portHex, 16, 32)
		if err != nil {
			continue
		}

		entries = append(entries, PortEntry{
			Protocol:     "tcp",
			LocalAddress: parts[0],
			LocalPort:    int(port),
			State:        "LISTEN",
		})
	}

	return entries, scanner.Err()
}
