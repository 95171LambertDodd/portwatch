package alerting

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// AlertType represents the category of a port alert.
type AlertType string

const (
	AlertNewBinding    AlertType = "NEW_BINDING"
	AlertConflict      AlertType = "CONFLICT"
	AlertUnexpected    AlertType = "UNEXPECTED"
	AlertRemovedBinding AlertType = "REMOVED_BINDING"
)

// Alert holds details about a detected port event.
type Alert struct {
	Timestamp time.Time
	Type      AlertType
	Entry     portscanner.PortEntry
	Message   string
}

// Alerter watches for changes in port state and emits alerts.
type Alerter struct {
	known  map[string]portscanner.PortEntry
	output io.Writer
}

// NewAlerter creates an Alerter that writes alerts to the given writer.
// If w is nil, os.Stdout is used.
func NewAlerter(w io.Writer) *Alerter {
	if w == nil {
		w = os.Stdout
	}
	return &Alerter{
		known:  make(map[string]portscanner.PortEntry),
		output: w,
	}
}

// Diff compares a fresh snapshot against the previously known state and
// returns alerts for newly appeared or disappeared bindings.
func (a *Alerter) Diff(current []portscanner.PortEntry) []Alert {
	var alerts []Alert

	currentMap := make(map[string]portscanner.PortEntry, len(current))
	for _, e := range current {
		key := entryKey(e)
		currentMap[key] = e

		if _, seen := a.known[key]; !seen {
			alerts = append(alerts, Alert{
				Timestamp: time.Now(),
				Type:      AlertNewBinding,
				Entry:     e,
				Message:   fmt.Sprintf("new binding detected on %s:%d (pid %d)", e.LocalAddr, e.LocalPort, e.PID),
			})
		}
	}

	// Detect bindings that were present before but are no longer in the snapshot.
	for key, e := range a.known {
		if _, stillPresent := currentMap[key]; !stillPresent {
			alerts = append(alerts, Alert{
				Timestamp: time.Now(),
				Type:      AlertRemovedBinding,
				Entry:     e,
				Message:   fmt.Sprintf("binding removed on %s:%d (pid %d)", e.LocalAddr, e.LocalPort, e.PID),
			})
		}
	}

	a.known = currentMap
	return alerts
}

// Emit writes a human-readable representation of the alert to the output writer.
func (a *Alerter) Emit(alert Alert) {
	fmt.Fprintf(a.output, "[%s] [%s] %s\n",
		alert.Timestamp.Format(time.RFC3339),
		alert.Type,
		alert.Message,
	)
}

func entryKey(e portscanner.PortEntry) string {
	return fmt.Sprintf("%s:%d:%s", e.LocalAddr, e.LocalPort, e.Protocol)
}
