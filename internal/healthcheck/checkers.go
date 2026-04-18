package healthcheck

import (
	"fmt"
	"os"
	"time"
)

// FileWritableChecker returns a Checker that verifies a path's parent dir is writable.
func FileWritableChecker(name, path string) Checker {
	return func() ComponentHealth {
		ch := ComponentHealth{Name: name, CheckedAt: time.Now()}
		f, err := os.CreateTemp(path, ".healthcheck-*")
		if err != nil {
			ch.Status = StatusDegraded
			ch.Message = fmt.Sprintf("cannot write to %s: %v", path, err)
			return ch
		}
		f.Close()
		os.Remove(f.Name())
		ch.Status = StatusOK
		return ch
	}
}

// StaticChecker always returns the given status. Useful for testing or stubs.
func StaticChecker(name string, status Status, msg string) Checker {
	return func() ComponentHealth {
		return ComponentHealth{
			Name:      name,
			Status:    status,
			Message:   msg,
			CheckedAt: time.Now(),
		}
	}
}
