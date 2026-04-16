package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	ScanInterval  time.Duration `yaml:"scan_interval"`
	AlertOutput   string        `yaml:"alert_output"`
	WatchedPorts  []int         `yaml:"watched_ports"`
	IgnoredPorts  []int         `yaml:"ignored_ports"`
	ProcNetTCP    string        `yaml:"proc_net_tcp"`
	ProcNetTCP6   string        `yaml:"proc_net_tcp6"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanInterval: 5 * time.Second,
		AlertOutput:  "stdout",
		WatchedPorts: []int{},
		IgnoredPorts: []int{},
		ProcNetTCP:   "/proc/net/tcp",
		ProcNetTCP6:  "/proc/net/tcp6",
	}
}

// Load reads a YAML config file from path and merges it with defaults.
// If path is empty, the default config is returned.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parsing YAML: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are semantically valid.
func (c *Config) Validate() error {
	if c.ScanInterval <= 0 {
		return fmt.Errorf("config: scan_interval must be positive")
	}
	if c.ProcNetTCP == "" {
		return fmt.Errorf("config: proc_net_tcp path must not be empty")
	}
	if c.ProcNetTCP6 == "" {
		return fmt.Errorf("config: proc_net_tcp6 path must not be empty")
	}
	return nil
}

// IsIgnored reports whether port is in the ignored list.
func (c *Config) IsIgnored(port int) bool {
	for _, p := range c.IgnoredPorts {
		if p == port {
			return true
		}
	}
	return false
}

// IsWatched reports whether port is explicitly watched.
// If WatchedPorts is empty, all ports are considered watched.
func (c *Config) IsWatched(port int) bool {
	if len(c.WatchedPorts) == 0 {
		return true
	}
	for _, p := range c.WatchedPorts {
		if p == port {
			return true
		}
	}
	return false
}

// ShouldAlert reports whether an alert should be raised for the given port.
// A port should trigger an alert if it is watched and not ignored.
func (c *Config) ShouldAlert(port int) bool {
	return c.IsWatched(port) && !c.IsIgnored(port)
}
