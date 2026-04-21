// Package scorecard computes a risk score for a port binding event
// based on configurable weighted criteria such as unexpected port, unknown process,
// or baseline violation.
package scorecard

import (
	"fmt"

	"github.com/user/portwatch/internal/portscanner"
)

// Rule defines a scoring criterion applied to a port entry.
type Rule struct {
	// Name is a human-readable label for the rule.
	Name string
	// Weight is added to the total score when the rule matches.
	Weight int
	// Match returns true when this rule applies to the given entry.
	Match func(e portscanner.Entry) bool
}

// Result holds the computed score and the names of rules that fired.
type Result struct {
	Score    int
	Reasons  []string
}

// Scorer evaluates port entries against a set of weighted rules.
type Scorer struct {
	rules []Rule
}

// New creates a Scorer from the provided rules.
// Returns an error if any rule has an empty name, nil Match, or negative weight.
func New(rules []Rule) (*Scorer, error) {
	for _, r := range rules {
		if r.Name == "" {
			return nil, fmt.Errorf("scorecard: rule name must not be empty")
		}
		if r.Match == nil {
			return nil, fmt.Errorf("scorecard: rule %q has nil Match function", r.Name)
		}
		if r.Weight < 0 {
			return nil, fmt.Errorf("scorecard: rule %q has negative weight %d", r.Name, r.Weight)
		}
	}
	return &Scorer{rules: rules}, nil
}

// Evaluate applies all rules to the entry and returns the aggregated Result.
func (s *Scorer) Evaluate(e portscanner.Entry) Result {
	var res Result
	for _, r := range s.rules {
		if r.Match(e) {
			res.Score += r.Weight
			res.Reasons = append(res.Reasons, r.Name)
		}
	}
	return res
}
