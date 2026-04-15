package baseline

import (
	"fmt"
	"io"

	"github.com/user/portwatch/internal/portscanner"
)

// Violation describes a port binding that is not in the baseline.
type Violation struct {
	Entry portscanner.PortEntry
	Reason string
}

// String returns a human-readable description of the violation.
func (v Violation) String() string {
	return fmt.Sprintf("VIOLATION proto=%s addr=%s port=%d pid=%d reason=%s",
		v.Entry.Proto, v.Entry.LocalAddress, v.Entry.LocalPort, v.Entry.PID, v.Reason)
}

// Checker compares live port entries against a Baseline.
type Checker struct {
	baseline *Baseline
	out      io.Writer
}

// NewChecker creates a Checker using the provided Baseline and output writer.
func NewChecker(b *Baseline, out io.Writer) *Checker {
	return &Checker{baseline: b, out: out}
}

// Check scans entries and returns any that are not covered by the baseline.
func (c *Checker) Check(entries []portscanner.PortEntry) []Violation {
	var violations []Violation
	for _, e := range entries {
		if !c.baseline.Contains(e.Proto, e.LocalAddress, e.LocalPort) {
			v := Violation{
				Entry:  e,
				Reason: "not in baseline",
			}
			violations = append(violations, v)
			if c.out != nil {
				fmt.Fprintln(c.out, v.String())
			}
		}
	}
	return violations
}
