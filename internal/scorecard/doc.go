// Package scorecard provides a rule-based risk scoring engine for port binding
// events detected by portwatch.
//
// Each Rule carries a Name, an integer Weight, and a Match predicate. When
// Scorer.Evaluate is called with a portscanner.Entry the scorer iterates all
// rules, accumulates the Weight of every matching rule, and returns a Result
// containing the total Score and the slice of Reason strings for matched rules.
//
// Usage:
//
//	rules := scorecard.DefaultRules()
//	rules = append(rules, scorecard.Rule{
//		Name:   "known-malware-port",
//		Weight: 100,
//		Match:  func(e portscanner.Entry) bool { return e.Port == 31337 },
//	})
//
//	s, err := scorecard.New(rules)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	res := s.Evaluate(entry)
//	fmt.Printf("risk score: %d reasons: %v\n", res.Score, res.Reasons)
package scorecard
