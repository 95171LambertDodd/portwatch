package portscanner

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single listening port binding observed on the host.
type PortEntry struct {
	Protocol     string `json:"protocol"`
	LocalAddress string `json:"local_address"`
	LocalPort    int    `json:"local_port"`
	PID          int    `json:"pid"`
	State        string `json:"state"`
}

// Scanner reads port bindings from /proc/net.
type Scanner struct {
	procNetTCP  string
	procNetTCP6 string
}

// NewScanner returns a Scanner pointed at the standard /proc/net paths.
func NewScanner() *Scanner {
	return &Scanner{
		procNetTCP:  "/proc/net/tcp",
		procNetTCP6: "/proc/net/tcp6",
	}
}

// Scan returns all currently listening TCP port entries.
func (s *Scanner) Scan() ([]PortEntry, error) {
	var entries []PortEntry
	for _, path := range []string{s.procNetTCP, s.procNetTCP6} {
		e, err := parseProcNet(path, "tcp")
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		entries = append(entries, e...)
	}
	return entries, nil
}

// parseProcNet parses a /proc/net/tcp[6] file and returns LISTEN entries.
func parseProcNet(path, proto string) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []PortEntry
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		state := fields[3]
		if state != "0A" { // 0A = TCP_LISTEN
			continue
		}
		addr, port, err := parseHexAddr(fields[1])
		if err != nil {
			continue
		}
		entries = append(entries, PortEntry{
			Protocol:     proto,
			LocalAddress: addr,
			LocalPort:    port,
			State:        "LISTEN",
		})
	}
	return entries, scanner.Err()
}

func parseHexAddr(s string) (string, int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid addr %q", s)
	}
	rawIP, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", 0, err
	}
	// reverse byte order for little-endian representation
	for i, j := 0, len(rawIP)-1; i < j; i, j = i+1, j-1 {
		rawIP[i], rawIP[j] = rawIP[j], rawIP[i]
	}
	ip := net.IP(rawIP).String()
	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return "", 0, err
	}
	return ip, int(port), nil
}
