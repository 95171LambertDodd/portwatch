package correlation

import (
	"encoding/json"
	"fmt"
	"os"
)

// RuleFile is the JSON structure for a correlation rules file.
type RuleFile struct {
	Rules []RuleEntry `json:"rules"`
}

// RuleEntry is the JSON-serialisable form of a Rule.
type RuleEntry struct {
	Port        uint16 `json:"port"`
	Protocol    string `json:"protocol"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Expected    bool   `json:"expected"`
}

// LoadFile reads a JSON rules file from path and returns a Correlator.
// If path is empty, an empty Correlator is returned.
func LoadFile(path string) (*Correlator, error) {
	if path == "" {
		return New(nil)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("correlation: open %s: %w", path, err)
	}
	defer f.Close()

	var rf RuleFile
	if err := json.NewDecoder(f).Decode(&rf); err != nil {
		return nil, fmt.Errorf("correlation: decode %s: %w", path, err)
	}

	rules := make([]Rule, 0, len(rf.Rules))
	for _, re := range rf.Rules {
		rules = append(rules, Rule{
			Port:     re.Port,
			Protocol: re.Protocol,
			Service: ServiceInfo{
				Name:        re.Name,
				Description: re.Description,
				Expected:    re.Expected,
			},
		})
	}
	return New(rules)
}
